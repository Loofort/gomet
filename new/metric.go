package metric

import (
	"fmt"
	"time"
)

// if bufSize is not set the defSize is used
var defSize = 1000

// default size for collector.In channel
// can be set by linker dirrective -X like:
// go build -ldflags "-X metric.bufSize Size"`
var bufSize string

type ErrOverload struct {
	error
}
type ErrLag ErrOverload

// default collector
var col = newCollector()

// collector collects all inbount event,
// when client ask it aggregates event with specified period and send to output
type collector struct {
	In  chan Event
	Out chan Stat
	Err chan error
}

// start collector process
// size specifies input channel length, if it is small the program could wait on sending metrics, and some time metrics could be distorted.
// in this case collector sends error.
// if size specified as 0 the default will be used
func newCollector(size int) collector {
	if size == 0 {
		size, _ = strconv.Atoi(bufSize)
		if size == 0 {
			// if no custom size iz set use default
			size = defSize
		}
	}

	in := make(chan Event, size)
	errc := make(chan error)

	go func() {

		border := size - 1
		events := make([]Event, 0, size)

		var ticker *time.Ticker

		// the borders for time period
		var low, high time.Time

		for {
			select {
			case evn := <-in:
				// if channel is full throw the error
				if len(in) >= border {
					go func() {
						err := ErrOverload{fmt.Errorf(`metrics are slow down the program. Current buffer size %d is not enought.`+
							`Please increase it, for default collector use compile derective: go build -ldflags "-X metric.bufSize Size"`, size)}
						errc <- err
					}()
					continue
				}

				events = append(events, evn)

			case t := <-ticker:
				// time to agreggate events, but there could be some pending events in In channel
				// this simple implementation doesn't deal with it, but raise an error.

				// walk throught events array and calculate Ticks
				var tick Tick
				for _, env := range events {
					// if new time period is started create new tick and old send to output
					if env.Time > high {
						go func(tick Tick) { out <- tick }(tick)
						low = high
						high = high + period
						tick = newTick(groups)
					}

					// throw error if event belongs to previous time period
					if env.Time <= low {
						go func() {
							err := ErrLag{fmt.Errorf(`event belong to previous time period. timeperiod=%d, lag=%d. try to increase period`, period, low.Sub(env.Time))}
							errc <- err
						}()
					}

					// agreggate metrics
					tick.Add(env)
				}
			}
		}
	}()
}

// Tick is the set of groups
type Tick map[string]Group

// group of similar metric contexts, e.g. each worker has own metric context but they are all belong to one group
type Group struct {
	Name    string
	Threads []Thread
}

// Thread represents state of metric context
type Thread struct {
	Start  time.Time
	States map[string]time.Duration
}

// Event e.g. state changing
type Event struct {
	Time time.Time
}

// ############ Metric type ##############
type Metric struct {
}

//func Goroutine(name string) Metric {}
func New(name string) Context {

}

func (m Metric) Close() {}

func (m Metric) State(state string) {

}

// integer up to 53bit precise can be safely converted to float64
func (m Metric) Value(val float64) {

}

// ############ Metrics for Channels #############
func ChanInc(name string) {}
func ChanDec(name string) {}
