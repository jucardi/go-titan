package prometheus

import (
	"context"
	"os"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var once sync.Once

var (
	metricClientInstance = &PrometheusClient{}
)

func GetSingleton() *PrometheusClient {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	once.Do(func() {
		metricClientInstance = &PrometheusClient{
			ctx:      context.Background(),
			hostname: hostname,
			env:      os.Getenv("EXECUTION_ENV"),
		}
	})
	return metricClientInstance
}

// ==== Start: Metric client ====

var RequestTimeGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: string(GaugeMetricNameRequestTime),
	Help: "Request time in milliseconds",
}, []string{"endpointName", "hostname", "env"})

var RequestTimeHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    string(HistogramMetricNameRequestTime),
	Help:    "Request time in milliseconds",
	Buckets: []float64{50, 100, 200, 300, 500, 800, 1300, 2100, 3400, 5500},
}, []string{"endpointName", "hostname", "env"})

var ApiErrorCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: string(CounterApiErrors),
	Help: "Total number of errors in API",
}, []string{"endpointName", "hostname", "env", "statusCode"})

func init() {
	prometheus.MustRegister(RequestTimeGauge)
	prometheus.MustRegister(RequestTimeHistogram)
	prometheus.MustRegister(ApiErrorCounter)
}

// ==== End: Metric client ====

type PrometheusClient struct {
	ctx      context.Context
	mode     string
	hostname string
	env      string
}

func (c *PrometheusClient) WithCtx(ctx context.Context) *PrometheusClient {
	return &PrometheusClient{
		ctx: ctx,
	}
}

func addDefaultLabels(labels map[string]string, c *PrometheusClient) map[string]string {
	labels["hostname"] = c.hostname
	labels["env"] = c.env
	return labels
}

func (c *PrometheusClient) SetGaugeValue(metricName GaugeMetricName, value float64, labels map[string]string) {
	labels = addDefaultLabels(labels, c)
	switch metricName {
	case GaugeMetricNameRequestTime:
		RequestTimeGauge.With(labels).Set(value)
	default:
		break
	}
}

func (c *PrometheusClient) ObserveHistogramValue(metricName HistogramMetricName, value float64, labels map[string]string) {
	labels = addDefaultLabels(labels, c)
	switch metricName {
	case HistogramMetricNameRequestTime:
		RequestTimeHistogram.With(labels).Observe(value)
	default:
		break
	}
}

func (c *PrometheusClient) IncreaseCounter(metricName CounterMetricName, labels map[string]string) {
	labels = addDefaultLabels(labels, c)
	switch metricName {
	case CounterApiErrors:
		ApiErrorCounter.With(labels).Inc()
	default:
		break
	}
}
