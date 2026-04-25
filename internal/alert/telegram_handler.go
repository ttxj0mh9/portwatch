package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const telegramAPIBase = "https://api.telegram.org/bot"

// TelegramHandler sends alert notifications via Telegram Bot API.
type TelegramHandler struct {
	token  string
	chatID string
	client *http.Client
}

// NewTelegramHandler creates a new TelegramHandler.
// token is the Telegram bot token and chatID is the target chat identifier.
func NewTelegramHandler(token, chatID string) (*TelegramHandler, error) {
	if token == "" {
		return nil, fmt.Errorf("telegram: bot token is required")
	}
	if chatID == "" {
		return nil, fmt.Errorf("telegram: chat ID is required")
	}
	return &TelegramHandler{
		token:  token,
		chatID: chatID,
		client: &http.Client{},
	}, nil
}

// Send dispatches an Event to the configured Telegram chat.
func (h *TelegramHandler) Send(e Event) error {
	text := FormatAlert(e)

	payload := map[string]string{
		"chat_id":    h.chatID,
		"text":       text,
		"parse_mode": "Markdown",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("telegram: failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s%s/sendMessage", telegramAPIBase, h.token)
	resp, err := h.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telegram: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram: unexpected status %d", resp.StatusCode)
	}
	return nil
}
