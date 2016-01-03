package gomet

import (
	"testing"
	"time"
)

// The test simulates use case represented by following shema:
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
// Scale = 2
// load of state1 = 50%
// load of state2 = 50%
// state1 avg complet time = 0 , beacouse incomplete
// state2 avg complet time = 0,5 tick
func Test1period(t *testing.T) {
	period := time.Second
	out := Setup(period)

	New("g").State("st1")
	m2 := New("g")
	m2.State("st2")
	time.Sleep(period / 2)
	m2.Close()
	New("g").State("st2")

	tick := <-out
	g := tick.Groups["g"]

	equalf(t, g.Scale(period), 2.0)
	equalf(t, g.Load("st1"), 0.5)
	equalf(t, g.Load("st2"), 0.5)
	equald(t, g.Lasted("st1"), 0)
	circad(t, g.Lasted("st2"), period/2)
}

// |*+|+ |.
// | *|*+|.
func Test2period(t *testing.T) {
	period := time.Second
	out := Setup(period)

	// send event
	go func() {
		m1 := New("g")
		m1.State("st1")
		time.Sleep(period / 2)
		m1.State("st2")
		m2 := New("g")
		m2.State("st1")
		time.Sleep(period)
		m1.Close()
		m2.State("st2")
	}()

	//wait two ticks
	tick := <-out
	t.Logf("tick1 : %#v", tick)
	equalf(t, tick.Scale("g"), 1.5)
	equalf(t, tick.Load("g", "st1"), 2.0/3)
	equalf(t, tick.Load("g", "st2"), 1.0/3)
	circad(t, tick.Lasted("g", "st1"), period/2)
	circad(t, tick.Lasted("g", "st2"), 0)

	tick = <-out
	t.Logf("tick2 : %#v", tick)
	equalf(t, tick.Scale("g"), 1.5)
	equalf(t, tick.Load("g", "st1"), 1.0/3)
	equalf(t, tick.Load("g", "st2"), 2.0/3)
	circad(t, tick.Lasted("g", "st1"), period)
	circad(t, tick.Lasted("g", "st2"), period)
}
