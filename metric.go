// metric allow you to collect metrics anywhere in the code,
// currently it supports convinient method for time and count metrics
// ussage:
//   m = metric.Now()
//   payload()
//   m.Set("Process Time")
//
// metrics are preserved in the map and wait till you grab them:
//   metrics := metric.Grab()
//   for _, tm := range metrics["Process Time"] {
//	   fmt.Println("time spend: ", tm.Seconds())
//   }
package metric

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var bufSize string
var store map[string][]int64
var soreErr error
var packChan chan pack
var mut sync.RWMutex

// init create a background goroutine that wait for metrics and puts it into global map
func init() {

	// bufSize can be set by linker dirrective -X
	// by default it is 1000
	size, _ := strconv.Atoi(bufSize)
	if size == 0 {
		size = 1000
	}
	border = size - 1

	packChan = make(chan pack, size)
	store = make(map[string][]int64, 10)

	go func() {
		for pack := range packChan {
			// lock if we Grab metrics
			mut.RLock()

			// if channel is full it could slow down the program payload processing
			// we have to erase buffer and mark the error
			if len(packChan) >= border {
				soreErr = fmt.Errorf(`metrics are slow down the program. Current buffer size %d is not enought. 
					Please increase it using compile derective: go build -ldflags "-X metric.bufSize Size"`, size)
			} else {

				// fill the map
				slice, ok := store[pack.Name]
				if !ok {
					slice = make([]int64, 0, size)
				}
				slice = append(slice, pack.Value)
				store[pack.Name] = slice
			}

			mut.RUnlock()
		}
	}()
}

// Grab returns all collected metrics.
// and erases internal collection.
// instead metrics it also can return error, that indicated that not all metrics was collected
func Grab() (map[string][]int64, error) {
	mut.Lock()
	defer mut.Unlock()

	metrics = store

	// erase internal store
	store = make(map[string][]int64, len(store))
	if soreErr != nil {
		return nil, err
	}

	return metrics, nil
}

type pack struct {
	Name  string
	Value int64
}

type Metric interface {
	Set(name string) Metric
}

type clock time.Time

// start new clock metric
func Now() Metric {
	m := clock(time.Now())
	return m
}

// Set calculates amount of time since clock was started and sends it to global metric's storage with given name
// It also return new clock metric that has been just started
func (c clock) Set(name string) Metric {
	tm := time.Now()
	dur := tm.Sub(time.Time(c))
	packChan <- pack{name, int64(dur)}
	return clock(tm)
}

type count int64

// create new count metric
func Count(cnt int) Metric {
	return count(cnt)
}

func (c count) Set(name string) Metric {
	packChan <- pack{name, int64(c)}
	return c
}
