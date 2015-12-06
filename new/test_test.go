// +build not

package gomet

import (
	"testing"
	"time"
)

// Test is simulating use case represented by  following shema:
// |**|.
// |+ |
// | +|.
//
// each line is goroutine (of the same group)
// | delimeter for ticks, so 1 tick consists of 2 poins
// * is state1
// + is state2
// . some unknown state, means continuing
// white space means no worker running
//
// expected the folowing stats:
// count of workers per tick = 2
// load of state1 = 50%
// load of state2 = 50%
// state1 avg complet time = 0 , beacouse incomplete
// state2 avg complet time = 0,5 tick
func TestBasic(t *testing.T) {
	in = make(chan Event, 1)
	out := collector(in, time.Microsecond)

	New("test").State("st1")
	m2 := New("test").State("st2")
	time.Sleep(500 * time.Nanosecond)
	m2.Close()
	New("test").State("st2")

	tick := <-out
	g := tick.Groups["test"]
	st1 := g.States["st1"]
	st2 := g.States["st2"]

	equalInt(g.Count, 2, t)
	equalFloat(st1.Load, 50, t)
	equalFloat(st2.Load, 50, t)
	equalDur(st1.Time, 0, t)
	equalDur(st2.Time, 500*time.Nanosecond, t)
}
