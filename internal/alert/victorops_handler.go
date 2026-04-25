package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// VictorOpsHandler sends alerts to VictorOps (Splunk On-Call) via the REST endpoint.
type VictorOpsHandler struct {
	restURL    string
	routingKey string
	client     *http.Client
}

type victorOpsPayload struct {
	MessageType       string `json:"message_type"`
	EntityID          string `json:"entity_id"`
	EntityDisplayName string `json:"entity_display_name"`
	StateMessage      string `json:"state_message"`
	Timestamp         int64  `json:"timestamp"`
}

// NewVictorOpsHandler creates a VictorOpsHandler.
// restURL is the base REST endpoint (e.g. https://alert.victorops.com/integrations/generic/...).
// routingKey determines which escalation policy receives the alert.
func NewVictorOpsHandler(restURL, routingKey string) (*VictorOpsHandler, error) {
	if restURL == "" {
		return nil, fmt.Errorf("victorops: rest_url is required")
	}
	if routingKey == "" {
		return nil, fmt.Errorf("victorops: routing_key is required")
	}
	return &VictorOpsHandler{
		restURL:    restURL,
		routingKey: routingKey,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send dispatches an Event to VictorOps.
func (v *VictorOpsHandler) Send(e Event) error {
	msgType := "INFO"
	if e.Level == LevelAlert {
		msgType = "CRITICAL"
	}

	payload := victorOpsPayload{
		MessageType:       msgType,
		EntityID:          fmt.Sprintf("portwatch-%s-%d", e.Change.Type, e.Change.Port),
		EntityDisplayName: fmt.Sprintf("portwatch: port %d %s", e.Change.Port, e.Change.Type),
		StateMessage:      FormatAlert(e.Change),
		Timestamp:         e.Time.Unix(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("victorops: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/%s", v.restURL, v.routingKey)
	resp, err := v.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("victorops: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("victorops: unexpected status %d", resp.StatusCode)
	}
	return nil
}
