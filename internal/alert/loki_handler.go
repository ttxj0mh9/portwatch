package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// lokiStream represents a single Loki log stream with labels and log entries.
type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][2]string       `json:"values"`
}

type lokiPushPayload struct {
	Streams []lokiStream `json:"streams"`
}

// LokiHandler sends alert events to a Grafana Loki instance via the push API.
type LokiHandler struct {
	pushURL string
	labels  map[string]string
	client  *http.Client
}

// NewLokiHandler creates a new LokiHandler.
// pushURL should be the full Loki push endpoint, e.g. http://localhost:3100/loki/api/v1/push.
func NewLokiHandler(pushURL string, labels map[string]string) (*LokiHandler, error) {
	if pushURL == "" {
		return nil, fmt.Errorf("loki: push URL is required")
	}
	if labels == nil {
		labels = map[string]string{"app": "portwatch"}
	}
	return &LokiHandler{
		pushURL: pushURL,
		labels:  labels,
		client:  &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send dispatches an Event to Loki as a structured log line.
func (h *LokiHandler) Send(ev Event) error {
	timestampNs := fmt.Sprintf("%d", ev.Time.UnixNano())
	msg := FormatAlert(ev)

	payload := lokiPushPayload{
		Streams: []lokiStream{
			{
				Stream: h.labels,
				Values: [][2]string{{timestampNs, msg}},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("loki: marshal payload: %w", err)
	}

	resp, err := h.client.Post(h.pushURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("loki: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("loki: unexpected status %d", resp.StatusCode)
	}
	return nil
}
