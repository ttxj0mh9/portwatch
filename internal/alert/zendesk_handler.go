package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/config"
)

// zendeskTicket represents the payload sent to the Zendesk Tickets API.
type zendeskTicket struct {
	Ticket zendeskTicketBody `json:"ticket"`
}

type zendeskTicketBody struct {
	Subject  string          `json:"subject"`
	Comment  zendeskComment  `json:"comment"`
	Priority string          `json:"priority"`
	Tags     []string        `json:"tags"`
}

type zendeskComment struct {
	Body string `json:"body"`
}

// ZendeskHandler sends alert events as Zendesk tickets.
type ZendeskHandler struct {
	subdomain string
	email     string
	apiToken  string
	client    *http.Client
}

// NewZendeskHandler creates a new ZendeskHandler from config.
func NewZendeskHandler(cfg config.ZendeskConfig) (*ZendeskHandler, error) {
	if cfg.Subdomain == "" {
		return nil, fmt.Errorf("zendesk: subdomain is required")
	}
	if cfg.Email == "" {
		return nil, fmt.Errorf("zendesk: email is required")
	}
	if cfg.APIToken == "" {
		return nil, fmt.Errorf("zendesk: api_token is required")
	}
	return &ZendeskHandler{
		subdomain: cfg.Subdomain,
		email:     cfg.Email,
		apiToken:  cfg.APIToken,
		client:    &http.Client{},
	}, nil
}

// Send creates a Zendesk ticket for the given event.
func (h *ZendeskHandler) Send(event Event) error {
	priority := "normal"
	if event.Level == LevelAlert {
		priority = "high"
	}

	tags := []string{"portwatch"}
	if event.Level == LevelAlert {
		tags = append(tags, "alert")
	} else {
		tags = append(tags, "info")
	}

	payload := zendeskTicket{
		Ticket: zendeskTicketBody{
			Subject:  fmt.Sprintf("[portwatch] %s", event.Summary),
			Comment:  zendeskComment{Body: FormatAlert(event)},
			Priority: priority,
			Tags:     tags,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("zendesk: marshal payload: %w", err)
	}

	url := fmt.Sprintf("https://%s.zendesk.com/api/v2/tickets.json", h.subdomain)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("zendesk: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(h.email+"/token", h.apiToken)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("zendesk: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("zendesk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
