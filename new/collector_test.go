package gomet

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func _equal(t *testing.T, got, expected interface{}) {
	if !reflect.DeepEqual(expected, got) {
		_, path, line, _ := runtime.Caller(2)
		_, file := filepath.Split(path)
		t.Fatalf("\n%s:%d got %#v (%s) but expected %#v (%s)\n", file, line, got, reflect.TypeOf(got), expected, reflect.TypeOf(expected))
	}
}
func equal(t *testing.T, got, expected interface{}) {
	// have to implement this wrapper in order to Caller func wroks properly
	_equal(t, got, expected)
}
func equald(t *testing.T, got, expected time.Duration) {
	_equal(t, got, expected)
}
func equalf(t *testing.T, got, expected float32) {
	var epsilon float32 = 0.0001
	if (got-expected) < epsilon && (expected-got) < epsilon {
		return
	}

	_equal(t, got, expected)
}

// send 10 events, expect receiving 10 events
func TestCollect(t *testing.T) {
	in := make(chan Event, 100)
	period := time.Millisecond
	eventsc := collect(in, period)

	go func() {
		for i := 0; i < 10; i++ {
			in <- Event{}
			time.Sleep(period / 10)
		}
		// have to sleep some time to allow collector process all events
		time.Sleep(period)
		close(in)
	}()

	cnt := 0
	for evs := range eventsc {
		// substract one auxiliary event
		evs = evs[:len(evs)-1]
		for range evs {
			cnt++
		}
		t.Logf("events in tick: %d", len(evs))
	}
	equal(t, cnt, 10)
}

// Test case [1 group, 1 tick]:
// |**|.
// |+ |
// | +|.
//
// expected:
// group scale = 2
// st1 load = 50%
// st2 load = 50%
// st1 time = 0 (beacouse incomplete)
// st2 time = 0,5 sec
func TestAggregate(t *testing.T) {
	period := time.Second
	half := 500 * time.Millisecond

	in := make(chan []Event)
	out := aggregate(in, period)

	start := time.Now()
	evs := []Event{
		Event{Group: "g", Worker: 1, State: "st1", Time: start},
		Event{Group: "g", Worker: 2, State: "st2", Time: start},
		Event{Group: "g", Worker: 2, State: "", Time: start.Add(half)},
		Event{Group: "g", Worker: 3, State: "st2", Time: start.Add(half)},
		Event{Group: "_auxiliary_", Time: start.Add(period)},
	}

	in <- evs
	tick := <-out
	close(in)
	g := tick.Groups["g"]

	equalf(t, g.Scale(period), 2)
	equalf(t, g.Load("st1"), 0.5)
	equalf(t, g.Load("st2"), 0.5)
	equald(t, g.Lasted("st1"), 0)
	equald(t, g.Lasted("st2"), 500*time.Millisecond)
}
