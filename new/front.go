// This package helps to gather metrics for goroutine execution time.
// The stats is grouped by specified unit of time.
// This package is simple, it doesn't send metrics across the network to some collector demon that aggregates and saves stats to DB.
// Instead the collector works in program space and provide you with stats "tick" by GO channel. You should save ticks by your own.
// This package is not intended to be used in huge enterprise soft, but demonstrates the concept of convinient metrics for golang programs.
package gomet

// Meter reports metric of particular goroutines to the backside collector
// Each goroutine should create own meter
type Meter struct {
}

// New creates Meter object. Meter is intend to report metrics of one goroutine/component.
// Several Meters could be created with the same name, they will be organized in group in stats.
// For example goroutine Workers should create Meters with identical name, one meter per Worker.
// It's good practics to always close meter in defer upter creation, like:
//   m := gomet.NewMeter("some.worker")
//   defer m.Close()
func NewMeter(name string) Meter {

}

// Close reports that particular goroutine's life circle is over
// it is safe to close meter twice, at second time it just does nothing.
func (m Meter) Close() {

}

// State sends to stats collector new goroutine state.
// Colector calcualte the average time duration that goroutine spends on each state
func (m Meter) State(state string) {

}

// ChanIn and ChanOut used together to measure time that object spends in channel
// Call ChanIn before put obj to chan, and ChanOut after get obj from chan.
// name should be unique for each channel.
func ChanIn(name string) {

}

// see ChanIn.
// ChanOut panics if ChanIn haven't been executed for the given name.
func ChanOut(name string) {

}
