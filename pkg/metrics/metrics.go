// pkg/metrics/metrics.go
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	MessagesRead = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "websocket_messages_read_total",
			Help: "Total number of messages read from the WebSocket.",
		},
	)
	ReadErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "websocket_read_errors_total",
			Help: "Total number of errors encountered while reading from the WebSocket.",
		},
	)
)

func init() {
	prometheus.MustRegister(MessagesRead)
	prometheus.MustRegister(ReadErrors)
}

// StartMetricsServer starts an HTTP server for Prometheus metrics.
func StartMetricsServer(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(addr, nil)
}
