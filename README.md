# Metrics

metric allow you to collect metrics anywhere in the code,
currently it supports convinient method for time and count metrics

## ussage:

	m = metric.Now()
	payload()
	m.Set("Process Time")

metrics are preserved in the gloabl map and wait till you grab them:

	metrics := metric.Grab()
	for _, tm := range metrics["Process Time"] {
		fmt.Printf("nanoseconds spend: %d", tm)
	}

