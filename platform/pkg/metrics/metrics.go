package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"service", "method", "code"})

	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
	}, []string{"service", "method", "code"})
)

func init() {
	for _, service := range []string{"order-service", "iam-service", "app"} {
		for _, method := range []string{"GET", "POST", "PUT", "DELETE"} {
			for _, code := range []string{"200", "400", "401", "404", "500"} {
				HTTPRequestsTotal.WithLabelValues(service, method, code)
			}
		}
	}
}

var httpServiceName = "app"

func SetServiceName(name string) {
	if name != "" {
		httpServiceName = name
	}
}

func ObserveHTTP(method string, status int, d time.Duration) {
	code := strconv.Itoa(status)
	HTTPRequestsTotal.WithLabelValues(httpServiceName, method, code).Inc()
	HTTPRequestDuration.WithLabelValues(httpServiceName, method, code).Observe(d.Seconds())
}
