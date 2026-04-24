package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookPayload is the JSON body sent to the webhook endpoint.
type WebhookPayload struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Port      int    `json:"port"`
	Event     string `json:"event"`
	Message   string `json:"message"`
}

// WebhookHandler sends alert events to an HTTP endpoint.
type WebhookHandler struct {
	url    string
	client *http.Client
}

// NewWebhookHandler creates a WebhookHandler that posts to the given URL.
func NewWebhookHandler(url string, timeout time.Duration) *WebhookHandler {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &WebhookHandler{
		url: url,
		client: &http.Client{Timeout: timeout},
	}
}

// Send encodes the event as JSON and POSTs it to the configured URL.
func (w *WebhookHandler) Send(e Event) error {
	payload := WebhookPayload{
		Timestamp: e.Time.UTC().Format(time.RFC3339),
		Level:     e.Level.String(),
		Port:      e.Port,
		Event:     string(e.Kind),
		Message:   fmt.Sprintf("port %d %s", e.Port, e.Kind),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post to %s: %w", w.url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d from %s", resp.StatusCode, w.url)
	}
	return nil
}
