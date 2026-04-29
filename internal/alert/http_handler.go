package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HTTPHandler sends alerts to a generic HTTP endpoint with configurable
// method, headers, and a JSON body template.
type HTTPHandler struct {
	url     string
	method  string
	headers map[string]string
	client  *http.Client
}

type httpPayload struct {
	Event     string    `json:"event"`
	Port      int       `json:"port"`
	Proto     string    `json:"proto"`
	Level     string    `json:"level"`
	Timestamp time.Time `json:"timestamp"`
}

// NewHTTPHandler creates an HTTPHandler. method defaults to POST if empty.
func NewHTTPHandler(url, method string, headers map[string]string) (*HTTPHandler, error) {
	if url == "" {
		return nil, fmt.Errorf("http handler: url is required")
	}
	if method == "" {
		method = http.MethodPost
	}
	return &HTTPHandler{
		url:     url,
		method:  method,
		headers: headers,
		client:  &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (h *HTTPHandler) Send(e Event) error {
	payload := httpPayload{
		Event:     e.Kind.String(),
		Port:      e.Port,
		Proto:     "tcp",
		Level:     e.Level.String(),
		Timestamp: e.At,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("http handler: marshal payload: %w", err)
	}
	req, err := http.NewRequest(h.method, h.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("http handler: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range h.headers {
		req.Header.Set(k, v)
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("http handler: send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("http handler: unexpected status %d", resp.StatusCode)
	}
	return nil
}
