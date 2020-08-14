package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	statsd "github.com/etsy/statsd/examples/go"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var requestDurationHistogram = func() *prometheus.HistogramVec {
	hv := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_milliseconds",
			Help:    "Incoming request duration in milliseconds",
			Buckets: []float64{0.5, 0.9, 0.99},
		},
		[]string{"app", "endpoint", "code", "method"},
	)

	prometheus.MustRegister(hv)
	return hv
}()

// Log middleware логирует входящие запросы, а также шлет метрики
func Log(next http.Handler,
	appName string,
	statsdClient *statsd.StatsdClient) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := &logResponseWriter{rw: rw}
		next.ServeHTTP(lrw, r)

		code := lrw.getStatusCode()

		fields := log.Fields{
			"url":    r.URL.String(),
			"d":      time.Since(start).String(),
			"status": code,
			"method": r.Method,
		}
		log.WithFields(fields).Info("request")

		metric := fmt.Sprintf("service.%s.request.%s.%d", appName, pathToMetric(r.URL.Path), code)
		d := time.Since(start).Milliseconds()
		statsdClient.Timing(metric, d)
		requestDurationHistogram.With(prometheus.Labels{
			"app":      appName,
			"endpoint": pathToMetric(r.URL.Path),
			"code":     fmt.Sprintf("%d", code),
			"method":   r.Method,
		}).Observe(float64(d))
	})
}

func pathToMetric(path string) string {
	return strings.Replace(strings.Trim(path, "/"), "/", "_", -1)
}

type logResponseWriter struct {
	rw         http.ResponseWriter
	statusCode int
}

func (l *logResponseWriter) Header() http.Header {
	return l.rw.Header()
}

func (l *logResponseWriter) Write(b []byte) (int, error) {
	return l.rw.Write(b)
}

func (l *logResponseWriter) WriteHeader(statusCode int) {
	l.statusCode = statusCode
	l.rw.WriteHeader(statusCode)
}

func (l *logResponseWriter) getStatusCode() int {
	if l.statusCode > 0 {
		return l.statusCode
	}
	return http.StatusOK
}
