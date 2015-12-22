// This package helps to gather metrics for goroutine execution time.
// The stats is grouped by specified unit of time.
// This package is simple, it doesn't send metrics across the network to some collector demon that aggregates and saves stats to DB.
// Instead the collector works in program space and provide you with stats "tick" by GO channel. You should save ticks by your own.
// This package is not intended to be used in huge enterprise soft, but demonstrates the concept of convinient metrics for golang programs.
package gomet

import (
	"sync/atomic"
	"time"
)

// gwid contains auto-incremented worker id for newly created meter
//
// lets assume that we will generate new worker(or chan send) every nanoseconds.
// It is 1 billion chan writing per second!
// in this case the wid will overflows in more than 500 years
var gwid int64

// dummy var
var tm time.Time

// Meter reports metric of particular goroutines to the backside collector
// Each goroutine should create own meter
type Meter struct {
	// provided group name
	Group string
	// auto-incremented worker id
	Wid int64
}

// New creates Meter object. Meter is intend to report metrics of one goroutine/component.
// Several Meters could be created with the same name, they will be organized in group in stats.
// For example goroutine Workers should create Meters with identical name, one meter per Worker.
// It's good practices to always close meter in defer after creation, like:
//   m := gomet.New("some.worker")
//   defer m.Close()
func New(name string) Meter {
	m := Meter{name, atomic.AddInt64(&gwid, 1)}
	c <- Event{m.Group, m.Wid, name, tm}
	return m
}

// Close reports that particular goroutine's life circle is over
// it is safe to close meter twice, at second time it just does nothing.
func (m Meter) Close() {
	c <- Event{m.Group, m.Wid, "", tm}
	m.Group = ""
}

// State sends to stats collector new goroutine state.
// Colector calculate the average time duration that goroutine spends on each state
func (m Meter) State(state string) {
	c <- Event{m.Group, m.Wid, state, tm}
}

// ChanIn and ChanOut used together to measure time that object spends in channel
// Call ChanIn before put obj to chan, and ChanOut after get obj from chan.
// chanName should be unique for each channel.
func ChanIn(chanName string) {
	c <- Event{chanName, 0, "chan", tm}
}

// see ChanIn.
// ChanOut panics if ChanIn haven't been executed for the given name.
func ChanOut(chanName string) {
	c <- Event{chanName, 0, "", tm}
}
