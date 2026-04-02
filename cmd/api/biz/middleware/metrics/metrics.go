package metrics

import (
	"bytes"
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
)

var (
	registerOnce sync.Once

	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "code"},
	)

	httpRequestDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "code"},
	)

	httpRequestsInFlight = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "api_http_requests_in_flight",
			Help: "Current number of in-flight HTTP requests.",
		},
		[]string{"method", "path"},
	)
)

func registerMetrics() {
	registerOnce.Do(func() {
		prometheus.MustRegister(httpRequestsTotal)
		prometheus.MustRegister(httpRequestDurationSeconds)
		prometheus.MustRegister(httpRequestsInFlight)
	})
}

func Middleware() app.HandlerFunc {
	registerMetrics()

	return func(ctx context.Context, c *app.RequestContext) {
		method := string(c.Method())
		path := string(c.Path())

		httpRequestsInFlight.WithLabelValues(method, path).Inc()
		start := time.Now()
		defer func() {
			httpRequestsInFlight.WithLabelValues(method, path).Dec()

			code := strconv.Itoa(c.Response.StatusCode())
			duration := time.Since(start).Seconds()

			httpRequestsTotal.WithLabelValues(method, path, code).Inc()
			httpRequestDurationSeconds.WithLabelValues(method, path, code).Observe(duration)
		}()

		c.Next(ctx)
	}
}

func Handler(ctx context.Context, c *app.RequestContext) {
	registerMetrics()

	mfs, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		c.String(500, "gather metrics failed: %v", err)
		return
	}

	var buf bytes.Buffer
	enc := expfmt.NewEncoder(&buf, expfmt.FmtText)
	for _, mf := range mfs {
		if err := enc.Encode(mf); err != nil {
			c.String(500, "encode metrics failed: %v", err)
			return
		}
	}

	c.Header("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	c.String(200, buf.String())
}
