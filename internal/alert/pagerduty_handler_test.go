package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewPagerDutyHandler_MissingKey(t *testing.T) {
	_, err := NewPagerDutyHandler("")
	if err == nil {
		t.Fatal("expected error for empty integration key")
	}
}

func TestNewPagerDutyHandler_Valid(t *testing.T) {
	h, err := NewPagerDutyHandler("test-key-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestPagerDutyHandler_Send_Success(t *testing.T) {
	var received pdPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	h, _ := NewPagerDutyHandler("key-abc")
	h.client.Transport = rewriteTransport(ts.URL)

	e := NewEvent(8080, EventOpened)
	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.RoutingKey != "key-abc" {
		t.Errorf("expected routing key 'key-abc', got %q", received.RoutingKey)
	}
	if received.EventAction != "trigger" {
		t.Errorf("expected event_action 'trigger', got %q", received.EventAction)
	}
	if received.Payload.Source != "portwatch" {
		t.Errorf("expected source 'portwatch', got %q", received.Payload.Source)
	}
}

func TestPagerDutyHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	h, _ := NewPagerDutyHandler("key-abc")
	h.client.Transport = rewriteTransport(ts.URL)

	e := NewEvent(9090, EventOpened)
	if err := h.Send(e); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestPagerDutyHandler_Send_AlertSeverity(t *testing.T) {
	var received pdPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	h, _ := NewPagerDutyHandler("key-abc")
	h.client.Transport = rewriteTransport(ts.URL)

	e := NewEvent(22, EventOpened)
	e.Level = LevelAlert
	h.Send(e)

	if received.Payload.Severity != "critical" {
		t.Errorf("expected severity 'critical' for LevelAlert, got %q", received.Payload.Severity)
	}
}
