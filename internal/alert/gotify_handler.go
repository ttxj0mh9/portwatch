package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GotifyHandler sends alerts to a self-hosted Gotify server.
type GotifyHandler struct {
	serverURL string
	token     string
	client    *http.Client
}

type gotifyPayload struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

// NewGotifyHandler creates a GotifyHandler. serverURL must not be empty and
// token must be a valid Gotify application token.
func NewGotifyHandler(serverURL, token string) (*GotifyHandler, error) {
	if strings.TrimSpace(serverURL) == "" {
		return nil, fmt.Errorf("gotify: server URL is required")
	}
	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("gotify: application token is required")
	}
	return &GotifyHandler{
		serverURL: strings.TrimRight(serverURL, "/"),
		token:     token,
		client:    &http.Client{},
	}, nil
}

func (g *GotifyHandler) Send(event Event) error {
	priority := 5
	if event.Level == LevelAlert {
		priority = 9
	}

	p := gotifyPayload{
		Title:    fmt.Sprintf("portwatch: port %d %s", event.Port, event.Kind),
		Message:  FormatAlert(event),
		Priority: priority,
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("gotify: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/message?token=%s", g.serverURL, g.token)
	resp, err := g.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("gotify: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("gotify: unexpected status %d", resp.StatusCode)
	}
	return nil
}
