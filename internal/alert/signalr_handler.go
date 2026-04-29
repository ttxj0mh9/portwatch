package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SignalRHandler sends alerts to an Azure SignalR / custom HTTP hub endpoint.
type SignalRHandler struct {
	hubURL    string
	accessKey string
	client    *http.Client
}

type signalRPayload struct {
	Target    string            `json:"target"`
	Arguments []signalRArgument `json:"arguments"`
}

type signalRArgument struct {
	Level   string `json:"level"`
	Message string `json:"message"`
	Port    int    `json:"port"`
	Time    string `json:"time"`
}

// NewSignalRHandler creates a SignalRHandler. hubURL is the REST endpoint for
// broadcasting a message (e.g. https://<svc>.service.signalr.net/api/v1/hubs/<hub>).
func NewSignalRHandler(hubURL, accessKey string) (*SignalRHandler, error) {
	if hubURL == "" {
		return nil, fmt.Errorf("signalr: hub URL is required")
	}
	if accessKey == "" {
		return nil, fmt.Errorf("signalr: access key is required")
	}
	return &SignalRHandler{
		hubURL:    hubURL,
		accessKey: accessKey,
		client:    &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send dispatches the event to the configured SignalR hub endpoint.
func (h *SignalRHandler) Send(e Event) error {
	payload := signalRPayload{
		Target: "portwatch",
		Arguments: []signalRArgument{
			{
				Level:   e.Level.String(),
				Message: FormatAlert(e),
				Port:    e.Port,
				Time:    e.Time.Format(time.RFC3339),
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("signalr: marshal payload: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, h.hubURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("signalr: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.accessKey)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("signalr: send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("signalr: unexpected status %d", resp.StatusCode)
	}
	return nil
}
