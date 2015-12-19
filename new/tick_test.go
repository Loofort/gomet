package gomet

import (
	"testing"
	"time"
)

// test case:
// |**|.
func TestTick_1G1W1S(t *testing.T) {
	period := time.Second
	tick := newTick(time.Now().Add(period), period)
	tick.set("g", 1, "st", period)

	equald(t, tick.TookTime("g", "st"), 0)
	equalf(t, tick.Load("g", "st"), 100)
	equalf(t, tick.Scale("g"), 1)

	g := tick.Groups["g"]
	equald(t, g.TookTime("st", period), 0)
	equalf(t, g.Load("st", period), 100)
	equalf(t, g.Scale(period), 1)
}

// test case:
// |**|.
// |++|.
func TestTick_1G2W2S(t *testing.T) {
	period := time.Second
	tick := newTick(time.Now().Add(period), period)
	tick.set("g", 1, "st1", period)
	tick.set("g", 2, "st2", period)

	equald(t, tick.TookTime("g", "st1"), 0)
	equald(t, tick.TookTime("g", "st2"), 0)
	equalf(t, tick.Load("g", "st1"), 50)
	equalf(t, tick.Load("g", "st2"), 50)
	equalf(t, tick.Scale("g"), 2)
}

// test case:
// |**|.
// |++|.
// | *|.
func TestTick_1G3W2S(t *testing.T) {
	period := time.Second
	tick := newTick(time.Now().Add(period), period)
	tick.set("g", 1, "st1", period)
	tick.set("g", 2, "st2", period)
	tick.set("g", 3, "st1", 500*time.Millisecond)

	equald(t, tick.TookTime("g", "st1"), 0)
	equald(t, tick.TookTime("g", "st2"), 0)
	equalf(t, tick.Load("g", "st1"), (100.0+0.0+50.0)/3.0)
	equalf(t, tick.Load("g", "st2"), (0.0+100.0+0.0)/3.0)
	equalf(t, tick.Scale("g"), 2.5)
}
