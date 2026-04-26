package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewMSTeamsAdaptiveHandler_MissingURL(t *testing.T) {
	_, err := NewMSTeamsAdaptiveHandler("")
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestNewMSTeamsAdaptiveHandler_Valid(t *testing.T) {
	h, err := NewMSTeamsAdaptiveHandler("https://example.com/webhook")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestMSTeamsAdaptiveHandler_Send_Success(t *testing.T) {
	var received []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, _ := NewMSTeamsAdaptiveHandler(ts.URL)
	h.client = ts.Client()

	e := Event{
		Port:    8080,
		Message: "port opened",
		Level:   LevelInfo,
		Time:    time.Now(),
	}
	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload adaptiveCardPayload
	if err := json.Unmarshal(received, &payload); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}
	if payload.Type != "message" {
		t.Errorf("expected type 'message', got %q", payload.Type)
	}
	if len(payload.Attachments) == 0 {
		t.Fatal("expected at least one attachment")
	}
}

func TestMSTeamsAdaptiveHandler_Send_AlertColor(t *testing.T) {
	var received []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, _ := NewMSTeamsAdaptiveHandler(ts.URL)
	h.client = ts.Client()

	e := Event{Port: 22, Message: "unexpected port", Level: LevelAlert, Time: time.Now()}
	_ = h.Send(e)

	var payload adaptiveCardPayload
	_ = json.Unmarshal(received, &payload)
	header := payload.Attachments[0].Content.Body[0]
	if header.Color != "Attention" {
		t.Errorf("expected Attention color for alert, got %q", header.Color)
	}
}

func TestMSTeamsAdaptiveHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	h, _ := NewMSTeamsAdaptiveHandler(ts.URL)
	h.client = ts.Client()

	e := Event{Port: 443, Message: "test", Level: LevelInfo, Time: time.Now()}
	if err := h.Send(e); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestMSTeamsAdaptiveHandler_Send_UnreachableURL(t *testing.T) {
	h, _ := NewMSTeamsAdaptiveHandler("http://127.0.0.1:1")
	e := Event{Port: 80, Message: "test", Level: LevelInfo, Time: time.Now()}
	if err := h.Send(e); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
