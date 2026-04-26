package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// MSTeamsAdaptiveHandler sends alerts via Microsoft Teams Adaptive Card format (Power Automate / Workflows).
type MSTeamsAdaptiveHandler struct {
	webhookURL string
	client     *http.Client
}

type adaptiveCardPayload struct {
	Type        string           `json:"type"`
	Attachments []adaptiveAttach `json:"attachments"`
}

type adaptiveAttach struct {
	ContentType string          `json:"contentType"`
	Content     adaptiveContent `json:"content"`
}

type adaptiveContent struct {
	Schema  string           `json:"$schema"`
	Type    string           `json:"type"`
	Version string           `json:"version"`
	Body    []adaptiveElement `json:"body"`
}

type adaptiveElement struct {
	Type   string `json:"type"`
	Text   string `json:"text"`
	Weight string `json:"weight,omitempty"`
	Size   string `json:"size,omitempty"`
	Color  string `json:"color,omitempty"`
}

// NewMSTeamsAdaptiveHandler creates a new MSTeamsAdaptiveHandler.
func NewMSTeamsAdaptiveHandler(webhookURL string) (*MSTeamsAdaptiveHandler, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("msteams adaptive: webhook URL is required")
	}
	return &MSTeamsAdaptiveHandler{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send dispatches an alert event as an Adaptive Card to Microsoft Teams.
func (h *MSTeamsAdaptiveHandler) Send(e Event) error {
	color := "Good"
	if e.Level == LevelAlert {
		color = "Attention"
	}

	payload := adaptiveCardPayload{
		Type: "message",
		Attachments: []adaptiveAttach{
			{
				ContentType: "application/vnd.microsoft.card.adaptive",
				Content: adaptiveContent{
					Schema:  "http://adaptivecards.io/schemas/adaptive-card.json",
					Type:    "AdaptiveCard",
					Version: "1.4",
					Body: []adaptiveElement{
						{Type: "TextBlock", Text: "PortWatch Alert", Weight: "Bolder", Size: "Medium", Color: color},
						{Type: "TextBlock", Text: e.Message},
						{Type: "TextBlock", Text: fmt.Sprintf("Port: %d | %s", e.Port, e.Time.Format(time.RFC3339))},
					},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("msteams adaptive: marshal payload: %w", err)
	}

	resp, err := h.client.Post(h.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("msteams adaptive: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("msteams adaptive: unexpected status %d", resp.StatusCode)
	}
	return nil
}
