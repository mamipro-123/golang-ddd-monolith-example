package observability

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	httpRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being served",
		},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(httpRequestsInFlight)
}

// MetricsMiddleware collects HTTP metrics
func MetricsMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	httpRequestsInFlight.Inc()
	defer httpRequestsInFlight.Dec()

	err := c.Next()
	if err != nil {
		httpRequestsTotal.WithLabelValues(
			c.Method(),
			c.Path(),
			"500",
		).Inc()
		httpRequestDuration.WithLabelValues(
			c.Method(),
			c.Path(),
		).Observe(time.Since(start).Seconds())
		return err
	}

	httpRequestsTotal.WithLabelValues(
		c.Method(),
		c.Path(),
		string(rune(c.Response().StatusCode())),
	).Inc()
	httpRequestDuration.WithLabelValues(
		c.Method(),
		c.Path(),
	).Observe(time.Since(start).Seconds())

	return nil
}

// StartMetricsServer starts a Prometheus metrics server
func StartMetricsServer(addr string) error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(addr, nil)
} 