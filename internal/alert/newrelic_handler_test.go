package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewNewRelicHandler_MissingAccountID(t *testing.T) {
	_, err := NewNewRelicHandler("", "key")
	if err == nil {
		t.Fatal("expected error for missing account ID")
	}
}

func TestNewNewRelicHandler_MissingAPIKey(t *testing.T) {
	_, err := NewNewRelicHandler("12345", "")
	if err == nil {
		t.Fatal("expected error for missing API key")
	}
}

func TestNewNewRelicHandler_Valid(t *testing.T) {
	h, err := NewNewRelicHandler("12345", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestNewRelicHandler_Send_Success(t *testing.T) {
	var received newRelicEvent
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Insert-Key") != "testkey" {
			t.Errorf("missing or wrong X-Insert-Key header")
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, _ := NewNewRelicHandler("12345", "testkey")
	h.url = ts.URL

	e := NewEvent(8080, EventOpened, time.Unix(1700000000, 0))
	if err := h.Send(e); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if received.Port != 8080 {
		t.Errorf("expected port 8080, got %d", received.Port)
	}
	if received.EventType != "PortWatchEvent" {
		t.Errorf("unexpected event type: %s", received.EventType)
	}
	if received.ChangeType != string(EventOpened) {
		t.Errorf("unexpected change type: %s", received.ChangeType)
	}
}

func TestNewRelicHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	h, _ := NewNewRelicHandler("12345", "badkey")
	h.url = ts.URL

	e := NewEvent(443, EventClosed, time.Now())
	if err := h.Send(e); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestNewRelicHandler_Send_UnreachableURL(t *testing.T) {
	h, _ := NewNewRelicHandler("12345", "key")
	h.url = "http://127.0.0.1:0"

	e := NewEvent(80, EventOpened, time.Now())
	if err := h.Send(e); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
