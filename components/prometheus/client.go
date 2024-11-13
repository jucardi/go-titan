package prometheus

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricClientInstance = &PrometheusClient{}
)

func GetSingleton() *PrometheusClient {
	return metricClientInstance
}

// ==== Start: Metric client ====

var RequestTimeGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: string(GaugeMetricNameRequestTime),
	Help: "Request time in milliseconds",
}, []string{"endpointName"})

var RequestTimeHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    string(HistogramMetricNameRequestTime),
	Help:    "Request time in milliseconds",
	Buckets: []float64{50, 100, 200, 300, 500, 800, 1300, 2100, 3400, 5500},
}, []string{"endpointName"})

var ApiErrorCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: string(CounterApiErrors),
	Help: "Total number of errors in API",
}, []string{"endpointName", "statusCode"})

func init() {
	prometheus.MustRegister(RequestTimeGauge)
	prometheus.MustRegister(RequestTimeHistogram)
}

// ==== End: Metric client ====

type PrometheusClient struct {
	ctx context.Context
}

func (c *PrometheusClient) WithCtx(ctx context.Context) *PrometheusClient {
	return &PrometheusClient{
		ctx: ctx,
	}
}

func (c *PrometheusClient) SetGaugeValue(metricName GaugeMetricName, value float64, labels map[string]string) {
	switch metricName {
	case GaugeMetricNameRequestTime:
		RequestTimeGauge.With(labels).Set(value)
	default:
		break
	}
}

func (c *PrometheusClient) ObserveHistogramValue(metricName HistogramMetricName, value float64, labels map[string]string) {
	switch metricName {
	case HistogramMetricNameRequestTime:
		RequestTimeHistogram.With(labels).Observe(value)
	default:
		break
	}
}

func (c *PrometheusClient) IncreaseCounter(metricName CounterMetricName, labels map[string]string) {
	switch metricName {
	case CounterApiErrors:
		ApiErrorCounter.With(labels).Inc()
	default:
		break
	}
}
