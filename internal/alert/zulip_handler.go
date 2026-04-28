package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// ZulipHandler sends alerts to a Zulip stream via the Zulip REST API.
type ZulipHandler struct {
	baseURL  string
	email    string
	apiKey   string
	stream   string
	topic    string
	client   *http.Client
}

type zulipPayload struct {
	Type    string `json:"type"`
	To      string `json:"to"`
	Topic   string `json:"topic"`
	Content string `json:"content"`
}

// NewZulipHandler creates a new ZulipHandler.
// baseURL should be the root of your Zulip instance, e.g. "https://yourorg.zulipchat.com".
func NewZulipHandler(baseURL, email, apiKey, stream, topic string) (*ZulipHandler, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("zulip: baseURL is required")
	}
	if email == "" {
		return nil, fmt.Errorf("zulip: bot email is required")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("zulip: API key is required")
	}
	if stream == "" {
		return nil, fmt.Errorf("zulip: stream is required")
	}
	if topic == "" {
		topic = "portwatch alerts"
	}
	return &ZulipHandler{
		baseURL: baseURL,
		email:   email,
		apiKey:  apiKey,
		stream:  stream,
		topic:   topic,
		client:  &http.Client{},
	}, nil
}

// Send delivers an Event to the configured Zulip stream.
func (h *ZulipHandler) Send(e Event) error {
	body := zulipPayload{
		Type:    "stream",
		To:      h.stream,
		Topic:   h.topic,
		Content: FormatAlert(e),
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("zulip: marshal payload: %w", err)
	}

	url := h.baseURL + "/api/v1/messages"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("zulip: create request: %w", err)
	}
	req.SetBasicAuth(h.email, h.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("zulip: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("zulip: unexpected status %d", resp.StatusCode)
	}
	return nil
}
