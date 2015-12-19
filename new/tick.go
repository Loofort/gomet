package gomet

import (
	"time"
)

// One tick represents the aggregated stats for specified time duration.
// The tick contains all goroutine groups and channels measured with metrics.
// The each tick's info could be saved to DB or file or printed on screen or whatever you need.
type Tick struct {
	Time   time.Time
	Period time.Duration
	Groups map[string]Group
}

func newTick(tm time.Time, period time.Duration) Tick {
	return Tick{
		Time:   tm,
		Period: period,
		Groups: make(map[string]Group),
	}
}

func (t Tick) set(group string, worker int64, state string, dur time.Duration) {

	g, ok := t.Groups[group]
	if !ok {
		g = newGroup(t.Period)
		t.Groups[group] = g
	}

	g.set(worker, state, dur)
}

// Group contains summarized stats for goroutines group or channel
// It also contains collection of all workers from group
type Group struct {
	Workers map[int64]Worker
	Period  time.Duration
}

func newGroup(period time.Duration) Group {
	return Group{
		Period:  period,
		Workers: make(map[int64]Worker),
	}
}

func (g Group) set(worker int64, state string, dur time.Duration) {
	w, ok := g.Workers[worker]
	if !ok {
		w = newWorker()
		g.Workers[worker] = w
	}

	w.set(state, dur)
}

// Sacle is the amount of workers per tick,
// e.g. if we have 2 workers, first is existing during whole tick and second only for half
// then scale will be 1+0.5 = 1.5
func (g Group) Scale() float32 {
	var dur time.Duration
	for _, w := range g.Workers {
		for _, s := range w.States {
			dur += s.Duration
		}
	}

	return float32(dur) / float32(g.Period)
}

// Load return average load of sate for group workers
// load means the percent of time spent for the state
func (g Group) Load(state string) float32 {
	var dur time.Duration
	for _, w := range g.Workers {
		s, ok := w.States[state]
		if !ok {
			continue
		}
		dur += s.Duration
	}

	load := 100 * float32(dur) / float32((time.Duration(len(g.Workers)) * g.Period))
	return load
}

// TookTime returns average time that sate needs to become completed
func (g Group) TookTime(state string) time.Duration {
	cnt := int64(0)
	var dur time.Duration
	for _, w := range g.Workers {
		s, ok := w.States[state]
		if !ok {
			continue
		}

		dur += s.Duration
		cnt += s.Count
		if s.Incomplete > 0 {
			dur -= s.Incomplete
			cnt--
		}
	}
	if cnt == 0 {
		return 0
	}

	return dur / time.Duration(cnt)
}

//------------------------------------------------
// Worker contains stats for one goroutine or one transfer throught channel
type Worker struct {
	States map[string]State
}

func newWorker() Worker {
	return Worker{
		States: make(map[string]State),
	}
}

func (w Worker) set(state string, dur time.Duration) {
	s := w.States[state]

	s.Count++
	s.Incomplete = dur
	s.Duration += dur
	if s.Min == 0 || s.Min > dur {
		s.Min += dur
	}

	if s.Max < dur {
		s.Max = dur
	}

	w.States[state] = s
}

// State represent time info such as
//   average time percent for the tick
type State struct {
	Count    int64
	Min, Max time.Duration
	Duration time.Duration
	// incomplete contains suration of state that is in progress by the end of the tick
	// if it is 0 it means that worker not in this state by the end of aggreagte period
	Incomplete time.Duration
}
