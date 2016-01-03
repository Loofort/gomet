# Metrics

gomet is the metrics collector designed to see what goroutines do.

## ussage:

at the beginning of program setup gomet specifying stats period:

    tickCh := gomet.Setup(time.Second) 

every second channel tickCh produces stats data, so you can listen on channel in some goroutines

	for tick := range tickCh {
		saveToCSV(tick)
		// or store in some other places
	}

Tick contains stats for each groups, see godoc for details.

At the begining of the goroutine create a new Meter, at the end of goroutine close the Meter:

	m := New("NetWorker")
	defer m.Close()

Than Specify states of goroutine:

	m.State("wait on chan")
	//data := <- dataCh
	m.State("payload")
	// payload(data)

To get stats how long data spend on channel use the ChanIn() and ChanOut() methods (see godoc)

