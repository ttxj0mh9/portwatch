package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SlackHandler sends alert notifications to a Slack webhook URL.
type SlackHandler struct {
	webhookURL string
	client     *http.Client
}

type slackPayload struct {
	Text string `json:"text"`
}

// NewSlackHandler creates a SlackHandler that posts to the given Slack
// incoming-webhook URL.
func NewSlackHandler(webhookURL string) *SlackHandler {
	return &SlackHandler{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send dispatches the event to Slack as a formatted message.
func (s *SlackHandler) Send(e Event) error {
	msg := FormatAlert(e)
	payload := slackPayload{Text: msg}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}
	return nil
}
