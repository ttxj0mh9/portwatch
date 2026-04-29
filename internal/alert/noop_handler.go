package alert

import "log"

// NoopHandler is a handler that discards all alerts.
// It is useful as a fallback or in testing when no real destination is configured.
type NoopHandler struct {
	logger *log.Logger
	verbose bool
}

// NoopOption configures a NoopHandler.
type NoopOption func(*NoopHandler)

// WithNoopVerbose enables debug logging when an event is discarded.
func WithNoopVerbose(v bool) NoopOption {
	return func(h *NoopHandler) {
		h.verbose = v
	}
}

// WithNoopLogger sets a custom logger on the handler.
func WithNoopLogger(l *log.Logger) NoopOption {
	return func(h *NoopHandler) {
		h.logger = l
	}
}

// NewNoopHandler returns a NoopHandler that silently drops every event.
func NewNoopHandler(opts ...NoopOption) *NoopHandler {
	h := &NoopHandler{
		logger: log.Default(),
	}
	for _, o := range opts {
		o(h)
	}
	return h
}

// Send discards the event. If verbose mode is enabled it logs the discarded event.
func (h *NoopHandler) Send(e Event) error {
	if h.verbose {
		h.logger.Printf("[noop] discarding event: port=%d action=%s level=%s",
			e.Port, e.Action, e.Level)
	}
	return nil
}
