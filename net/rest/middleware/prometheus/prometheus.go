package prometheus

import (
	"time"

	"github.com/jucardi/go-titan/components/prometheus"
	"github.com/jucardi/go-titan/net/rest"
)

func Handler(c *rest.Context) {
	start := time.Now()
	c.Next()
	elapsed := time.Since(start)
	elapsedMs := float64(elapsed.Milliseconds())
	endpointName := c.Request.Method + " " + c.Request.URL.Path

	client := prometheus.GetSingleton()

	// Observe the elapsed time
	labels := make(map[string]string)
	labels["endpointName"] = endpointName
	client.ObserveHistogramValue(prometheus.HistogramMetricNameRequestTime, elapsedMs, labels)
	client.SetGaugeValue(prometheus.GaugeMetricNameRequestTime, elapsedMs, labels)

	// Increase error counter
	status := c.Writer.Status()
	if status >= 400 {
		counterLabels := make(map[string]string)
		counterLabels["statusCode"] = string(rune(status))
		counterLabels["endpointName"] = endpointName
		client.IncreaseCounter(prometheus.CounterApiErrors, counterLabels)
	}
}
