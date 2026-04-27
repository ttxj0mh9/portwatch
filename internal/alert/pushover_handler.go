package alert

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const pushoverAPIURL = "https://api.pushover.net/1/messages.json"

// PushoverHandler sends alerts via the Pushover notification service.
type PushoverHandler struct {
	userKey  string
	apiToken string
	client   *http.Client
}

// NewPushoverHandler creates a new PushoverHandler.
// userKey is the Pushover user or group key.
// apiToken is the application API token.
func NewPushoverHandler(userKey, apiToken string) (*PushoverHandler, error) {
	if userKey == "" {
		return nil, fmt.Errorf("pushover: user key is required")
	}
	if apiToken == "" {
		return nil, fmt.Errorf("pushover: api token is required")
	}
	return &PushoverHandler{
		userKey:  userKey,
		apiToken: apiToken,
		client:   &http.Client{},
	}, nil
}

// Send delivers an alert event via Pushover.
func (h *PushoverHandler) Send(event Event) error {
	priority := 0
	if event.Level == LevelAlert {
		priority = 1
	}

	params := url.Values{}
	params.Set("token", h.apiToken)
	params.Set("user", h.userKey)
	params.Set("title", fmt.Sprintf("portwatch: port %d %s", event.Port, event.Change))
	params.Set("message", FormatAlert(event))
	params.Set("priority", fmt.Sprintf("%d", priority))

	resp, err := h.client.Post(pushoverAPIURL, "application/x-www-form-urlencoded",
		strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("pushover: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		return fmt.Errorf("pushover: unexpected status %d", resp.StatusCode)
	}
	return nil
}
