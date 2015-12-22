package gomet

import (
	"time"
)

// application contains current state of all program components
// it is a map of groups, each group is map of workers, each worker is in particular state that starts at some time point
type app map[string]map[int64]worker

func newApp() app {
	return make(app)
}

// update app state with new event,
// return previous sate and duration for it
func (a app) update(ev Event) (string, time.Duration) {
	// get group
	g, ok := a[ev.Group]
	if !ok {
		g = make(map[int64]worker)
		a[ev.Group] = g
	}

	w, ok := g[ev.Worker]
	g[ev.Worker] = worker{ev.State, ev.Time} // set new state
	if !ok {
		// this is worker's start, no further processing is needed
		return "", 0
	}

	if ev.State == "" {
		// this is wroker's end, remove it form app sate
		delete(g, ev.Worker)
	}

	// calculate duration for event and
	dur := ev.Time.Sub(w.Start)
	return w.State, dur
}

type worker struct {
	State string
	Start time.Time
}
