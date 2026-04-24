package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const opsGenieAPIURL = "https://api.opsgenie.com/v2/alerts"

// OpsGenieHandler sends alerts to OpsGenie.
type OpsGenieHandler struct {
	apiKey string
	client *http.Client
}

// NewOpsGenieHandler creates a new OpsGenieHandler.
// Returns an error if apiKey is empty.
func NewOpsGenieHandler(apiKey string) (*OpsGenieHandler, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("opsGenie: api key must not be empty")
	}
	return &OpsGenieHandler{
		apiKey: apiKey,
		client: &http.Client{},
	}, nil
}

type opsGeniePayload struct {
	Message     string `json:"message"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

// Send dispatches an alert event to OpsGenie.
func (h *OpsGenieHandler) Send(e Event) error {
	priority := "P3"
	if e.Level == LevelAlert {
		priority = "P1"
	}

	payload := opsGeniePayload{
		Message:     fmt.Sprintf("portwatch: %s", e.Message),
		Description: fmt.Sprintf("Port %s — %s", e.Port, e.Message),
		Priority:    priority,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("opsGenie: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, opsGenieAPIURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("opsGenie: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+h.apiKey)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("opsGenie: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("opsGenie: unexpected status %d", resp.StatusCode)
	}
	return nil
}
