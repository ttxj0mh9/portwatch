package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewAmplitudeHandler_MissingAPIKey(t *testing.T) {
	_, err := NewAmplitudeHandler("", "")
	if err == nil {
		t.Fatal("expected error for missing API key, got nil")
	}
}

func TestNewAmplitudeHandler_Valid(t *testing.T) {
	h, err := NewAmplitudeHandler("test-key", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.endpoint != amplitudeDefaultEndpoint {
		t.Errorf("expected default endpoint %q, got %q", amplitudeDefaultEndpoint, h.endpoint)
	}
}

func TestAmplitudeHandler_Send_Success(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, _ := NewAmplitudeHandler("key-123", ts.URL)
	e := Event{
		Type:    EventOpened,
		Port:    8080,
		Level:   LevelAlert,
		Message: "port 8080 opened",
		Time:    time.Now(),
	}

	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["api_key"] != "key-123" {
		t.Errorf("expected api_key 'key-123', got %v", received["api_key"])
	}
	events, ok := received["events"].([]interface{})
	if !ok || len(events) != 1 {
		t.Fatalf("expected 1 event, got %v", received["events"])
	}
	ev := events[0].(map[string]interface{})
	if ev["event_type"] != "port_opened" {
		t.Errorf("expected event_type 'port_opened', got %v", ev["event_type"])
	}
}

func TestAmplitudeHandler_Send_ClosedEvent(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, _ := NewAmplitudeHandler("key-abc", ts.URL)
	e := Event{Type: EventClosed, Port: 443, Level: LevelInfo, Time: time.Now()}

	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events := received["events"].([]interface{})
	ev := events[0].(map[string]interface{})
	if ev["event_type"] != "port_closed" {
		t.Errorf("expected event_type 'port_closed', got %v", ev["event_type"])
	}
}

func TestAmplitudeHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	h, _ := NewAmplitudeHandler("key-xyz", ts.URL)
	e := Event{Type: EventOpened, Port: 22, Level: LevelAlert, Time: time.Now()}

	if err := h.Send(e); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}
