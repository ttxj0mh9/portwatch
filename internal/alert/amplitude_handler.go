package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const amplitudeDefaultEndpoint = "https://api2.amplitude.com/2/httpapi"

// AmplitudeHandler sends port change events to Amplitude Analytics.
type AmplitudeHandler struct {
	apiKey   string
	endpoint string
	client   *http.Client
}

// NewAmplitudeHandler creates a new AmplitudeHandler.
// Returns an error if apiKey is empty.
func NewAmplitudeHandler(apiKey, endpoint string) (*AmplitudeHandler, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("amplitude: API key is required")
	}
	if endpoint == "" {
		endpoint = amplitudeDefaultEndpoint
	}
	return &AmplitudeHandler{
		apiKey:   apiKey,
		endpoint: endpoint,
		client:   &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send dispatches an Event to Amplitude as a track call.
func (h *AmplitudeHandler) Send(e Event) error {
	type amplitudeEvent struct {
		EventType  string                 `json:"event_type"`
		EventProps map[string]interface{} `json:"event_properties"`
		Time       int64                  `json:"time"`
	}
	type amplitudePayload struct {
		APIKey string           `json:"api_key"`
		Events []amplitudeEvent `json:"events"`
	}

	eventType := "port_opened"
	if e.Type == EventClosed {
		eventType = "port_closed"
	}

	payload := amplitudePayload{
		APIKey: h.apiKey,
		Events: []amplitudeEvent{
			{
				EventType: eventType,
				EventProps: map[string]interface{}{
					"port":    e.Port,
					"level":   e.Level.String(),
					"message": e.Message,
				},
				Time: e.Time.UnixMilli(),
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("amplitude: marshal payload: %w", err)
	}

	resp, err := h.client.Post(h.endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("amplitude: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("amplitude: unexpected status %d", resp.StatusCode)
	}
	return nil
}
