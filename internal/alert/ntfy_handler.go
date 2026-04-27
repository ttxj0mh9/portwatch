package alert

import (
	"fmt"
	"net/http"
	"strings"
)

// NtfyHandler sends alerts to a ntfy.sh topic (self-hosted or cloud).
type NtfyHandler struct {
	serverURL string
	topic     string
	client    *http.Client
}

// NewNtfyHandler creates a new NtfyHandler.
// serverURL is the base URL (e.g. "https://ntfy.sh" or a self-hosted instance).
// topic is the ntfy topic to publish to.
func NewNtfyHandler(serverURL, topic string) (*NtfyHandler, error) {
	if strings.TrimSpace(serverURL) == "" {
		return nil, fmt.Errorf("ntfy: server URL is required")
	}
	if strings.TrimSpace(topic) == "" {
		return nil, fmt.Errorf("ntfy: topic is required")
	}
	return &NtfyHandler{
		serverURL: strings.TrimRight(serverURL, "/"),
		topic:     topic,
		client:    &http.Client{},
	}, nil
}

// Send publishes an Event to the configured ntfy topic.
func (h *NtfyHandler) Send(e Event) error {
	url := fmt.Sprintf("%s/%s", h.serverURL, h.topic)

	body := strings.NewReader(FormatAlert(e))
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("ntfy: failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Title", fmt.Sprintf("portwatch: port %d %s", e.Port, e.Kind))

	if e.Level == LevelAlert {
		req.Header.Set("Priority", "high")
		req.Header.Set("Tags", "warning,rotating_light")
	} else {
		req.Header.Set("Priority", "default")
		req.Header.Set("Tags", "information_source")
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("ntfy: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy: unexpected status %d", resp.StatusCode)
	}
	return nil
}
