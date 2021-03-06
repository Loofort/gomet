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

func (t Tick) set(group string, wid int64, state string, start time.Time, dur time.Duration, inProgress bool) {

	g, ok := t.Groups[group]
	if !ok {
		g = make(Group)
	}

	w, ok := g[wid]
	if !ok {
		w = make(Worker)
	}

	s := w[state].next(dur)
	if inProgress {
		s.Incomplete = dur
	}
	if s.Count == 1 && t.Time.Sub(start) > t.Period {
		s.Age = t.Time.Sub(start) - t.Period
		s.Duration -= s.Age
	}

	w[state] = s
	g[wid] = w
	t.Groups[group] = g
}

// Scale is the amount of workers per tick,
// e.g. if we have 2 workers, first is existing during whole tick and second only for half
// then scale will be 1+0.5 = 1.5
// It is shorthand for tick.Groups[group].Scale(tick.Period)
func (t Tick) Scale(group string) float32 {
	return t.Groups[group].Scale(t.Period)
}

// Count returns number of times that group was in particular state during period
// It is shorthand for t.Groups[group].Count(state)
func (t Tick) Count(group, state string) int {
	return t.Groups[group].Count(state)
}

// Load returns average load of sate for workers group
// load means the percent of time spent for the state during tick period
// It is shorthand for tick.Groups[group].Load(state, tick.Period)
func (t Tick) Load(group, state string) float32 {
	return t.Groups[group].Load(state)
}

// Lasted  returns average time that state needs to get completed
// It is shorthand for tick.Groups[group].Lasted(state)
func (t Tick) Lasted(group, state string) time.Duration {
	return t.Groups[group].Lasted(state)
}

//---------------------------------------

// Group contains summarized stats for goroutines group or channel
// It also contains collection of all workers from group
type Group map[int64]Worker

// Sacle is the amount of workers per tick,
// e.g. if we have 2 workers, first is existing during whole tick and second only for half
// then scale will be 1+0.5 = 1.5
func (g Group) Scale(period time.Duration) float32 {
	return float32(g.lifeSum()) / float32(period)
}

// calculates sum of workers live time during period
func (g Group) lifeSum() (dur time.Duration) {
	for _, w := range g {
		for _, s := range w {
			dur += s.Duration
		}
	}
	return dur
}

// Count return count of times group been in particular state
func (g Group) Count(state string) int {
	cnt := int64(0)
	for _, w := range g {
		s, ok := w[state]
		if !ok {
			continue
		}
		cnt += s.Count
	}
	return int(cnt)
}

// Load returns average load of sate for workers group
// load means the percent of time spent for the state
func (g Group) Load(state string) float32 {
	if len(g) == 0 {
		return 0
	}

	var dur time.Duration
	for _, w := range g {
		s, ok := w[state]
		if !ok {
			continue
		}
		dur += s.Duration
	}

	lifeSum := g.lifeSum()
	if lifeSum == 0 {
		return 0
	}
	return float32(dur) / float32(lifeSum)
}

// Count returns number of times that group was in particular state during period
func (g Group) Lasted(state string) time.Duration {
	cnt := int64(0)
	var dur time.Duration
	for _, w := range g {
		s, ok := w[state]
		if !ok {
			continue
		}

		if s.Incomplete == 0 {
			dur += s.Duration + s.Age
			cnt += s.Count
		}
	}
	if cnt == 0 {
		return 0
	}

	return dur / time.Duration(cnt)
}

// Worker contains stats for one goroutine or one transfer through channel
type Worker map[string]State

// State represent time info such as
//   average time percent for the tick
type State struct {
	Count    int64
	Min, Max time.Duration
	// time spent in previous periods
	Age time.Duration
	// time spent since event start or period start
	Duration time.Duration

	// incomplete contains duration of state that is in progress by the end of the tick
	// if it is 0 it means that worker not in this state by the end of aggregate period
	Incomplete time.Duration
}

// next modify state stat and returns new state.
func (s State) next(dur time.Duration) State {
	s.Count++
	s.Duration += dur
	if s.Min == 0 || s.Min > dur {
		s.Min += dur
	}

	if s.Max < dur {
		s.Max = dur
	}

	return s
}
