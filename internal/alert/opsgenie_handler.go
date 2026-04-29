package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/config"
)

const opsgenieAPIURL = "https://api.opsgenie.com/v2/alerts"

type OpsGenieHandler struct {
	apiKey  string
	tags    []string
	client  *http.Client
	apiURL  string
}

func NewOpsGenieHandler(cfg config.OpsGenieConfig) (*OpsGenieHandler, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("opsgenie: api_key is required")
	}
	apiURL := cfg.APIURL
	if apiURL == "" {
		apiURL = opsgenieAPIURL
	}
	return &OpsGenieHandler{
		apiKey: cfg.APIKey,
		tags:   cfg.Tags,
		client: &http.Client{},
		apiURL: apiURL,
	}, nil
}

func (h *OpsGenieHandler) Send(event Event) error {
	priority := "P3"
	if event.Level == LevelAlert {
		priority = "P1"
	}

	payload := map[string]interface{}{
		"message":     event.Message,
		"description": fmt.Sprintf("Port %s — %s", event.Port, event.Message),
		"priority":    priority,
		"tags":        h.tags,
		"source":      "portwatch",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("opsgenie: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, h.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("opsgenie: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+h.apiKey)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("opsgenie: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("opsgenie: unexpected status %d", resp.StatusCode)
	}
	return nil
}
