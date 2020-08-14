package outgoing

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	clientNameLabel = "client"
	successLabel    = "success"

	requestHandled = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "client_requests_handled",
			Help: "Requests handled",
		},
		[]string{clientNameLabel, successLabel},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "client_requests_duration_seconds",
			Help: "Request duration",
		},
		[]string{clientNameLabel},
	)
)

// NewMetricsStorage creates and returns a new MetricsStorage.
func NewMetricsStorage(name string) *MetricsStorage {
	return &MetricsStorage{name: name}
}

// MetricsStorage reports the metrics to prometheus.
type MetricsStorage struct {
	name string
}

// ReportRequest reports request
func (r *MetricsStorage) ReportRequest(success bool) {
	labels := map[string]string{
		clientNameLabel: r.name,
	}
	if success {
		labels[successLabel] = "true"
	} else {
		labels[successLabel] = "false"
	}

	requestHandled.With(labels).Inc()
}

// ReportDuration reports request duration
func (r *MetricsStorage) ReportDuration(duration time.Duration) {
	labels := map[string]string{
		clientNameLabel: r.name,
	}
	requestDuration.With(labels).Observe(duration.Seconds())
}

// nolint:gochecknoinits (LOL)
func init() {
	prometheus.MustRegister(requestHandled)
	prometheus.MustRegister(requestDuration)
}

type metrics interface {
	ReportRequest(success bool)
	ReportDuration(duration time.Duration)
}

// MetricsMw describes metrics middleware struct
type MetricsMw struct {
	next    http.RoundTripper
	metrics metrics
}

// NewMetricsMw returns new metrics middleware
func NewMetricsMw(name string, next http.RoundTripper) *MetricsMw {
	return &MetricsMw{
		next:    next,
		metrics: NewMetricsStorage(name),
	}
}

// RoundTrip implements http.RoundTripper interface
func (red MetricsMw) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	res, err := red.next.RoundTrip(req)

	success := true
	if err != nil {
		success = false
	}
	red.metrics.ReportRequest(success)
	red.metrics.ReportDuration(time.Since(start))

	return res, err
}
