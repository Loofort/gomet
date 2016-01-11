package gomet

import (
	"sync/atomic"
	"time"
)

// global chan that used by meter to send events
var c = closedChan()

// it returns closed Event chan.
// So if you forgot to call Setup the Meter will panic on sending to closed chan
func closedChan() chan Event {
	c := make(chan Event)
	close(c)
	return c
}

// Setup starts internal metrics collector.
// The collector groups metrics with specified time period.
// I suggest not setting period less than a second.
// Setup should be executed only once at the program start.
func Setup(period time.Duration) chan Tick {
	// 1024 is big enough to smooth activity spike.
	// but if profiler shows permanent awaiting on c chan the collector should be added with one more goroutine reader from c
	c = make(chan Event, 1000)
	out := collector(c, period)
	return out
}

func collector(in chan Event, period time.Duration) chan Tick {
	// collect events
	eventsc := collect(in, period)

	// aggregate events
	tickc := aggregate(eventsc, period)
	return tickc
}

func collect(in chan Event, period time.Duration) chan []Event {
	out := make(chan []Event, 100)
	go func() {
		defer close(out)

		ticker := time.NewTicker(period)
		defer ticker.Stop()

		buf := make([]Event, 0, 100)
		wids := make(map[string]struct {
			top int64
			low int64
		})

		for {
			select {
			case ev, ok := <-in:
				if !ok {
					return
				}
				// to simplify implementation we neglect the time spending by event in input channel
				// todo: for the high resolution stats it could be vital and needs to be improved
				ev.Time = time.Now()

				// for chan events calculate worker id
				if ev.Worker == 0 {
					// we in chan event
					w := wids[ev.Group]

					wid := &w.top //ChanIn
					if ev.State == "" {
						wid = &w.low //ChanOut
					}

					ev.Worker = atomic.AddInt64(wid, 1)
					wids[ev.Group] = w
				}

				buf = append(buf, ev)

			case now := <-ticker.C:
				evs := make([]Event, len(buf), len(buf)+1)
				copy(evs, buf)
				buf = buf[:0]

				// add auxiliary event to transfer tick time
				evs = append(evs, Event{Group: "_auxiliary_", Time: now})
				out <- evs
			}
		}
	}()
	return out
}

func aggregate(in chan []Event, period time.Duration) chan Tick {
	out := make(chan Tick, 100)
	go func() {
		defer close(out)

		a := newApp()
		var tick Tick
		var outChan chan Tick
		queue := make([]Tick, 0, 1)

		for {
			select {
			case evs, ok := <-in:
				if !ok {
					return
				}
				// pull auxiliary event to get tick time
				aux := evs[len(evs)-1]
				evs = evs[:len(evs)-1]

				t := newTick(aux.Time, period)
				for _, ev := range evs {
					// get goroutine state and duration
					state, start, dur := a.update(ev)
					if state == "" {
						continue
					}

					// add stats to tick
					t.set(ev.Group, ev.Worker, state, start, dur, false)
				}

				// states that are runing - are in app, but not in tick
				// states that have finished - not in app, but are in tick
				for gname, g := range a {
					for wid, aw := range g {
						dur := t.Time.Sub(aw.Start)
						t.set(gname, wid, aw.State, aw.Start, dur, true)
					}
				}

				// initiate sending of tick to output channel if no pending tick
				if outChan == nil {
					tick = t
					outChan = out
					break
				}

				// put tick in the awaiting queue
				queue = append(queue, t)

			case outChan <- tick:
				if len(queue) == 0 {
					outChan = nil
					break
				}

				// get next tick from queue
				tick = queue[0]
				n := copy(queue, queue[1:])
				queue = queue[:n]
			}
		}
	}()
	return out
}
