package gomet

import (
	"testing"
	"time"
)

// test case:
// |**|.
func TestTick_1G1W1S(t *testing.T) {
	period := time.Second
	now := time.Now()
	tick := newTick(now.Add(period), period)
	tick.set("g", 1, "st", now, period, true)

	equald(t, tick.Lasted("g", "st"), 0)
	equalf(t, tick.Load("g", "st"), 1)
	equalf(t, tick.Scale("g"), 1)

	g := tick.Groups["g"]
	equald(t, g.Lasted("st"), 0)
	equalf(t, g.Load("st"), 1)
	equalf(t, g.Scale(period), 1)
}

// test case:
// |**|.
// |++|.
func TestTick_1G2W2S(t *testing.T) {
	period := time.Second
	now := time.Now()
	tick := newTick(now.Add(period), period)
	tick.set("g", 1, "st1", now, period, true)
	tick.set("g", 2, "st2", now, period, true)

	equald(t, tick.Lasted("g", "st1"), 0)
	equald(t, tick.Lasted("g", "st2"), 0)
	equalf(t, tick.Load("g", "st1"), 0.5)
	equalf(t, tick.Load("g", "st2"), 0.5)
	equalf(t, tick.Scale("g"), 2)
}

// test case:
// |**|.
// |++|.
// | *|.
func TestTick_1G3W2S(t *testing.T) {
	period := time.Second
	now := time.Now()
	tick := newTick(now.Add(period), period)
	tick.set("g", 1, "st1", now, period, true)
	tick.set("g", 2, "st2", now, period, true)
	tick.set("g", 3, "st1", now, 500*time.Millisecond, true)

	equald(t, tick.Lasted("g", "st1"), 0)
	equald(t, tick.Lasted("g", "st2"), 0)
	equalf(t, tick.Load("g", "st1"), (1.0+0.0+0.5)/(1.0+1.0+0.5))
	equalf(t, tick.Load("g", "st2"), (0.0+1.0+0.0)/(1.0+1.0+0.5))
	equalf(t, tick.Scale("g"), 2.5)
}
