package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// splunkEvent is the HEC (HTTP Event Collector) payload.
type splunkEvent struct {
	Time   float64        `json:"time"`
	Source string         `json:"source"`
	Event  map[string]any `json:"event"`
}

// SplunkHandler sends alerts to a Splunk HEC endpoint.
type SplunkHandler struct {
	hecURL string
	token  string
	source string
	client *http.Client
}

// NewSplunkHandler creates a SplunkHandler.
// hecURL must be the full HEC endpoint, e.g. https://splunk:8088/services/collector.
func NewSplunkHandler(hecURL, token, source string) (*SplunkHandler, error) {
	if hecURL == "" {
		return nil, fmt.Errorf("splunk: HEC URL is required")
	}
	if token == "" {
		return nil, fmt.Errorf("splunk: HEC token is required")
	}
	if source == "" {
		source = "portwatch"
	}
	return &SplunkHandler{
		hecURL: hecURL,
		token:  token,
		source: source,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send dispatches an Event to Splunk HEC.
func (h *SplunkHandler) Send(e Event) error {
	payload := splunkEvent{
		Time:   float64(e.Time.Unix()),
		Source: h.source,
		Event: map[string]any{
			"port":    e.Port,
			"action":  e.Action,
			"level":   e.Level.String(),
			"message": fmt.Sprintf("port %d %s", e.Port, e.Action),
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("splunk: marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, h.hecURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("splunk: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Splunk "+h.token)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("splunk: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("splunk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
