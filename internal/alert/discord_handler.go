package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// DiscordHandler sends alert notifications to a Discord channel via webhook.
type DiscordHandler struct {
	webhookURL string
	client     *http.Client
}

type discordPayload struct {
	Content string         `json:"content,omitempty"`
	Embeds  []discordEmbed `json:"embeds,omitempty"`
}

type discordEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color"`
}

// NewDiscordHandler creates a new DiscordHandler for the given webhook URL.
func NewDiscordHandler(webhookURL string) (*DiscordHandler, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("discord: webhook URL must not be empty")
	}
	return &DiscordHandler{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}, nil
}

// Send delivers an Event to the configured Discord webhook.
func (h *DiscordHandler) Send(e Event) error {
	color := 0x36a64f // green for info
	if e.Level == LevelAlert {
		color = 0xff0000 // red for alert
	}

	payload := discordPayload{
		Embeds: []discordEmbed{
			{
				Title:       fmt.Sprintf("[%s] Port %s", e.Level, e.Port),
				Description: e.Message,
				Color:       color,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("discord: marshal payload: %w", err)
	}

	resp, err := h.client.Post(h.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("discord: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord: unexpected status %d", resp.StatusCode)
	}
	return nil
}
