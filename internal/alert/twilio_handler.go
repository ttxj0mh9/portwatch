package alert

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/user/portwatch/internal/config"
)

// TwilioHandler sends SMS alerts via the Twilio REST API.
type TwilioHandler struct {
	accountSID string
	authToken  string
	from       string
	to         []string
	client     *http.Client
}

// NewTwilioHandler creates a TwilioHandler from config.
func NewTwilioHandler(cfg config.TwilioConfig) (*TwilioHandler, error) {
	if cfg.AccountSID == "" {
		return nil, fmt.Errorf("twilio: account_sid is required")
	}
	if cfg.AuthToken == "" {
		return nil, fmt.Errorf("twilio: auth_token is required")
	}
	if cfg.From == "" {
		return nil, fmt.Errorf("twilio: from_number is required")
	}
	if len(cfg.To) == 0 {
		return nil, fmt.Errorf("twilio: at least one to_number is required")
	}
	return &TwilioHandler{
		accountSID: cfg.AccountSID,
		authToken:  cfg.AuthToken,
		from:       cfg.From,
		to:         cfg.To,
		client:     &http.Client{},
	}, nil
}

// Send dispatches an SMS for each configured recipient.
func (h *TwilioHandler) Send(event Event) error {
	body := FormatAlert(event)
	apiURL := fmt.Sprintf(
		"https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json",
		h.accountSID,
	)
	for _, recipient := range h.to {
		form := url.Values{}
		form.Set("From", h.from)
		form.Set("To", recipient)
		form.Set("Body", body)

		req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(form.Encode()))
		if err != nil {
			return fmt.Errorf("twilio: build request: %w", err)
		}
		req.SetBasicAuth(h.accountSID, h.authToken)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := h.client.Do(req)
		if err != nil {
			return fmt.Errorf("twilio: send to %s: %w", recipient, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			var apiErr struct {
				Message string `json:"message"`
			}
			_ = json.NewDecoder(resp.Body).Decode(&apiErr)
			return fmt.Errorf("twilio: unexpected status %d for %s: %s", resp.StatusCode, recipient, apiErr.Message)
		}
	}
	return nil
}
