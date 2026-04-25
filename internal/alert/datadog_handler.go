package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const datadogEventsURL = "https://api.datadoghq.com/api/v1/events"

// DatadogHandler sends alert events to the Datadog Events API.
type DatadogHandler struct {
	apiKey string
	url    string
	client *http.Client
}

type datadogEvent struct {
	Title     string   `json:"title"`
	Text      string   `json:"text"`
	AlertType string   `json:"alert_type"`
	Tags      []string `json:"tags,omitempty"`
	SourceTypeName string `json:"source_type_name"`
}

// NewDatadogHandler creates a DatadogHandler. apiKey must not be empty.
func NewDatadogHandler(apiKey string) (*DatadogHandler, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("datadog: api key is required")
	}
	return &DatadogHandler{
		apiKey: apiKey,
		url:    datadogEventsURL,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send dispatches an Event to Datadog.
func (h *DatadogHandler) Send(e Event) error {
	alertType := "info"
	if e.Level == LevelAlert {
		alertType = "warning"
	}

	payload := datadogEvent{
		Title:          fmt.Sprintf("portwatch: port %d %s", e.Port, e.Change),
		Text:           FormatAlert(e),
		AlertType:      alertType,
		Tags:           []string{fmt.Sprintf("port:%d", e.Port), fmt.Sprintf("change:%s", e.Change)},
		SourceTypeName: "portwatch",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("datadog: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, h.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("datadog: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DD-API-KEY", h.apiKey)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("datadog: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("datadog: unexpected status %d", resp.StatusCode)
	}
	return nil
}
