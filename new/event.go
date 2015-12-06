package gomet

import (
	"time"
)

// Event is the stat object that Meter sends to Collector each time Meter's Method is called.
// Collector has representation of current meter state and changes it according to Event.
type Event struct {
	Group  string
	Worker int64
	State  string
	Time   time.Time
}
