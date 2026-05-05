package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const newRelicEventsURL = "https://insights-collector.newrelic.com/v1/accounts/%s/events"

// NewRelicHandler sends port change events to New Relic Insights.
type NewRelicHandler struct {
	accountID string
	apiKey    string
	url       string
	client    *http.Client
}

type newRelicEvent struct {
	EventType   string `json:"eventType"`
	Port        int    `json:"port"`
	ChangeType  string `json:"changeType"`
	Level       string `json:"level"`
	Message     string `json:"message"`
	Timestamp   int64  `json:"timestamp"`
}

// NewNewRelicHandler creates a New Relic Insights event handler.
func NewNewRelicHandler(accountID, apiKey string) (*NewRelicHandler, error) {
	if accountID == "" {
		return nil, fmt.Errorf("newrelic: account ID is required")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("newrelic: API key is required")
	}
	return &NewRelicHandler{
		accountID: accountID,
		apiKey:    apiKey,
		url:       fmt.Sprintf(newRelicEventsURL, accountID),
		client:    &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send dispatches an Event to New Relic as a custom insight event.
func (h *NewRelicHandler) Send(e Event) error {
	payload := newRelicEvent{
		EventType:  "PortWatchEvent",
		Port:       e.Port,
		ChangeType: string(e.Type),
		Level:      e.Level.String(),
		Message:    FormatAlert(e),
		Timestamp:  e.Time.Unix(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("newrelic: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, h.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("newrelic: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Insert-Key", h.apiKey)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("newrelic: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("newrelic: unexpected status %d", resp.StatusCode)
	}
	return nil
}
