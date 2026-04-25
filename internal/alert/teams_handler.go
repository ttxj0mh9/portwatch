package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// teamsPayload is the Adaptive Card message format for Microsoft Teams.
type teamsPayload struct {
	Type       string          `json:"type"`
	Attachments []teamsAttachment `json:"attachments"`
}

type teamsAttachment struct {
	ContentType string       `json:"contentType"`
	Content     teamsContent `json:"content"`
}

type teamsContent struct {
	Schema  string        `json:"$schema"`
	Type    string        `json:"type"`
	Version string        `json:"version"`
	Body    []teamsBlock  `json:"body"`
}

type teamsBlock struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Wrap  bool   `json:"wrap"`
	Style string `json:"style,omitempty"`
}

// TeamsHandler sends alerts to a Microsoft Teams channel via an incoming webhook.
type TeamsHandler struct {
	webhookURL string
	client     *http.Client
}

// NewTeamsHandler creates a TeamsHandler. webhookURL must be non-empty.
func NewTeamsHandler(webhookURL string) (*TeamsHandler, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("teams: webhook URL is required")
	}
	return &TeamsHandler{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}, nil
}

// Send dispatches an Event to the configured Teams webhook.
func (h *TeamsHandler) Send(e Event) error {
	body := teamsPayload{
		Type: "message",
		Attachments: []teamsAttachment{
			{
				ContentType: "application/vnd.microsoft.card.adaptive",
				Content: teamsContent{
					Schema:  "http://adaptivecards.io/schemas/adaptive-card.json",
					Type:    "AdaptiveCard",
					Version: "1.4",
					Body: []teamsBlock{
						{Type: "TextBlock", Text: fmt.Sprintf("**[%s] portwatch alert**", e.Level), Wrap: true, Style: "heading"},
						{Type: "TextBlock", Text: e.Message, Wrap: true},
						{Type: "TextBlock", Text: fmt.Sprintf("Port: %s | Time: %s", e.Port, e.Time.Format(timeFormat)), Wrap: true},
					},
				},
			},
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("teams: marshal payload: %w", err)
	}

	resp, err := h.client.Post(h.webhookURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("teams: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("teams: unexpected status %d", resp.StatusCode)
	}
	return nil
}
