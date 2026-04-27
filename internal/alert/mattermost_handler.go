package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/alert"
)

// mattermostPayload is the incoming webhook payload for Mattermost.
type mattermostPayload struct {
	Text     string `json:"text"`
	Username string `json:"username,omitempty"`
	IconURL  string `json:"icon_url,omitempty"`
}

// MattermostHandler sends alerts to a Mattermost incoming webhook.
type MattermostHandler struct {
	webhookURL string
	username   string
	iconURL    string
	client     *http.Client
}

// NewMattermostHandler creates a new MattermostHandler.
// webhookURL must be non-empty.
func NewMattermostHandler(webhookURL, username, iconURL string) (*MattermostHandler, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("mattermost: webhook URL is required")
	}
	return &MattermostHandler{
		webhookURL: webhookURL,
		username:   username,
		iconURL:    iconURL,
		client:     &http.Client{},
	}, nil
}

// Send delivers an alert event to the configured Mattermost webhook.
func (h *MattermostHandler) Send(evt Event) error {
	text := fmt.Sprintf("**[portwatch]** %s", FormatAlert(evt))

	payload := mattermostPayload{
		Text:     text,
		Username: h.username,
		IconURL:  h.iconURL,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("mattermost: marshal payload: %w", err)
	}

	resp, err := h.client.Post(h.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("mattermost: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mattermost: unexpected status %d", resp.StatusCode)
	}
	return nil
}
