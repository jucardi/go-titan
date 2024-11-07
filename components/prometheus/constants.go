package prometheus

type GaugeMetricName string

const (
	GaugeMetricNameRequestTime GaugeMetricName = "gauge_request_time_milliseconds"
)

type HistogramMetricName string

const (
	HistogramMetricNameRequestTime HistogramMetricName = "histogram_request_time_milliseconds"
)

type CounterMetricName string

const (
	CounterApiErrors CounterMetricName = "counter_api_errors"
)
