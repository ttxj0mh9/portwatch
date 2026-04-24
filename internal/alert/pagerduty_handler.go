package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const pagerDutyEventsURL = "https://events.pagerduty.com/v2/enqueue"

// PagerDutyHandler sends alerts to PagerDuty via the Events API v2.
type PagerDutyHandler struct {
	integrationKey string
	client         *http.Client
}

// NewPagerDutyHandler creates a new PagerDutyHandler with the given integration key.
func NewPagerDutyHandler(integrationKey string) (*PagerDutyHandler, error) {
	if integrationKey == "" {
		return nil, fmt.Errorf("pagerduty: integration key must not be empty")
	}
	return &PagerDutyHandler{
		integrationKey: integrationKey,
		client:         &http.Client{Timeout: 10 * time.Second},
	}, nil
}

type pdPayload struct {
	RoutingKey  string    `json:"routing_key"`
	EventAction string    `json:"event_action"`
	Payload     pdDetails `json:"payload"`
}

type pdDetails struct {
	Summary  string `json:"summary"`
	Source   string `json:"source"`
	Severity string `json:"severity"`
}

// Send dispatches an Event to PagerDuty.
func (h *PagerDutyHandler) Send(e Event) error {
	severity := "info"
	if e.Level == LevelAlert {
		severity = "critical"
	}

	body := pdPayload{
		RoutingKey:  h.integrationKey,
		EventAction: "trigger",
		Payload: pdDetails{
			Summary:  FormatAlert(e),
			Source:   "portwatch",
			Severity: severity,
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("pagerduty: failed to marshal payload: %w", err)
	}

	resp, err := h.client.Post(pagerDutyEventsURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("pagerduty: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
