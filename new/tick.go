package gomet

import (
	"time"
)

// One tick represents the aggregated stats for specified time duration.
// The tick contains all goroutine groups and channels measured with metrics.
// The each tick's info could be saved to DB or file or printed on screen or whatever you need.
type Tick struct {
	Time       time.Time
	Components map[string]Component
}

// Component contains summarized stats for goroutines group or channel
// It also contains collection of all workers from group
type Component struct {
	//  amount of workers per tick, e.g. if we got 2 worker, first was exists during whole tick and second only for half than total amount will be 1.5
	Amount  float32
	States  map[name]State
	Workers []Worker
}

// Worker contains stats for one goroutine or one transfer throught channel
type Worker struct {
	States map[name]State
}

// State represent time info such as
//   average time percent for the tick
type State struct {
}
