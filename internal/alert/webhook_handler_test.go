package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWebhookHandler_Send_Success(t *testing.T) {
	var received WebhookPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewWebhookHandler(ts.URL, time.Second)
	e := NewEvent(8080, EventOpened, fixedTime())

	if err := h.Send(e); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	if received.Port != 8080 {
		t.Errorf("expected port 8080, got %d", received.Port)
	}
	if received.Event != string(EventOpened) {
		t.Errorf("expected event %q, got %q", EventOpened, received.Event)
	}
	if received.Timestamp == "" {
		t.Error("expected non-empty timestamp")
	}
}

func TestWebhookHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewWebhookHandler(ts.URL, time.Second)
	e := NewEvent(443, EventClosed, fixedTime())

	if err := h.Send(e); err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestWebhookHandler_Send_UnreachableURL(t *testing.T) {
	h := NewWebhookHandler("http://127.0.0.1:1", 200*time.Millisecond)
	e := NewEvent(80, EventOpened, fixedTime())

	if err := h.Send(e); err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
