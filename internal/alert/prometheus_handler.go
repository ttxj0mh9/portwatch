package alert

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusHandler exposes port-change events as Prometheus metrics.
type PrometheusHandler struct {
	mu      sync.Mutex
	opened  prometheus.Counter
	closed  prometheus.Counter
	alerts  prometheus.Counter
	server  *http.Server
}

// NewPrometheusHandler registers metrics and starts an HTTP server on addr
// (e.g. "127.0.0.1:9090") at path (e.g. "/metrics").
func NewPrometheusHandler(addr, path string) (*PrometheusHandler, error) {
	if addr == "" {
		return nil, fmt.Errorf("prometheus handler: addr must not be empty")
	}
	if path == "" {
		path = "/metrics"
	}

	reg := prometheus.NewRegistry()

	opened := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "portwatch_ports_opened_total",
		Help: "Total number of ports observed opening.",
	})
	closed := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "portwatch_ports_closed_total",
		Help: "Total number of ports observed closing.",
	})
	alerts := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "portwatch_alerts_total",
		Help: "Total number of alert-level events dispatched.",
	})

	reg.MustRegister(opened, closed, alerts)

	mux := http.NewServeMux()
	mux.Handle(path, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	srv := &http.Server{Addr: addr, Handler: mux}
	go srv.ListenAndServe() //nolint:errcheck

	return &PrometheusHandler{
		opened: opened,
		closed: closed,
		alerts: alerts,
		server: srv,
	}, nil
}

// Send implements Handler. It increments the appropriate counters.
func (h *PrometheusHandler) Send(event Event) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	switch event.Kind {
	case KindOpened:
		h.opened.Inc()
	case KindClosed:
		h.closed.Inc()
	}

	if event.Level == LevelAlert {
		h.alerts.Inc()
	}
	return nil
}

// Close shuts down the embedded HTTP server.
func (h *PrometheusHandler) Close() error {
	return h.server.Close()
}
