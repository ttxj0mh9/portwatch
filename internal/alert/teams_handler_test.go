package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewTeamsHandler_MissingURL(t *testing.T) {
	_, err := NewTeamsHandler("")
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestNewTeamsHandler_Valid(t *testing.T) {
	h, err := NewTeamsHandler("https://example.webhook.office.com/webhook/abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestTeamsHandler_Send_Success(t *testing.T) {
	var received teamsPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, _ := NewTeamsHandler(ts.URL)
	h.client = ts.Client()

	e := Event{
		Level:   LevelAlert,
		Message: "port opened",
		Port:    "8080/tcp",
		Time:    time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}

	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Type != "message" {
		t.Errorf("expected type 'message', got %q", received.Type)
	}
	if len(received.Attachments) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(received.Attachments))
	}
	blocks := received.Attachments[0].Content.Body
	if len(blocks) != 3 {
		t.Fatalf("expected 3 body blocks, got %d", len(blocks))
	}
}

func TestTeamsHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	h, _ := NewTeamsHandler(ts.URL)
	h.client = ts.Client()

	err := h.Send(Event{Level: LevelInfo, Port: "443/tcp", Time: time.Now()})
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestTeamsHandler_Send_UnreachableURL(t *testing.T) {
	h, _ := NewTeamsHandler("http://127.0.0.1:19999/unreachable")
	err := h.Send(Event{Level: LevelAlert, Port: "22/tcp", Time: time.Now()})
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
