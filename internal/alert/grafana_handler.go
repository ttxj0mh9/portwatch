package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// grafanaPayload represents a Grafana annotation payload.
type grafanaPayload struct {
	Text    string   `json:"text"`
	Tags    []string `json:"tags"`
	Time    int64    `json:"time"`
}

// GrafanaHandler sends port change events as Grafana annotations.
type GrafanaHandler struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewGrafanaHandler creates a new GrafanaHandler.
// baseURL is the Grafana instance URL (e.g. http://localhost:3000).
// apiKey is a Grafana API key with Editor or Admin role.
func NewGrafanaHandler(baseURL, apiKey string) (*GrafanaHandler, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("grafana: baseURL is required")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("grafana: apiKey is required")
	}
	return &GrafanaHandler{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{},
	}, nil
}

// Send posts an annotation to Grafana for the given event.
func (h *GrafanaHandler) Send(ev Event) error {
	tags := []string{"portwatch"}
	if ev.Type == EventOpened {
		tags = append(tags, "opened")
	} else {
		tags = append(tags, "closed")
	}

	payload := grafanaPayload{
		Text: FormatAlert(ev),
		Tags: tags,
		Time: ev.Time.UnixMilli(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("grafana: marshal payload: %w", err)
	}

	url := h.baseURL + "/api/annotations"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("grafana: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.apiKey)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("grafana: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("grafana: unexpected status %d", resp.StatusCode)
	}
	return nil
}
