package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/snapshot"
)

// NewRocketChatHandler creates an alert handler that posts messages to a Rocket.Chat
// incoming webhook URL.
func NewRocketChatHandler(webhookURL, channel, username string) (*rocketChatHandler, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("rocketchat: webhook URL is required")
	}
	return &rocketChatHandler{
		webhookURL: webhookURL,
		channel:    channel,
		username:   username,
		client:     &http.Client{},
	}, nil
}

type rocketChatHandler struct {
	webhookURL string
	channel    string
	username   string
	client     *http.Client
}

type rocketChatPayload struct {
	Text     string `json:"text"`
	Channel  string `json:"channel,omitempty"`
	Username string `json:"username,omitempty"`
	Emoji    string `json:"emoji,omitempty"`
}

func (h *rocketChatHandler) Send(event Event, snap snapshot.Snapshot) error {
	emoji := ":information_source:"
	if event.Level == LevelAlert {
		emoji = ":warning:"
	}

	payload := rocketChatPayload{
		Text:     fmt.Sprintf("%s *Port %s*: %s", emoji, snap.Ports[0].String(), event.Message),
		Channel:  h.channel,
		Username: h.username,
		Emoji:    emoji,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("rocketchat: marshal payload: %w", err)
	}

	resp, err := h.client.Post(h.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("rocketchat: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("rocketchat: unexpected status %d", resp.StatusCode)
	}
	return nil
}
