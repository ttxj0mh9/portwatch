package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func TestNewVictorOpsHandler_MissingURL(t *testing.T) {
	_, err := NewVictorOpsHandler("", "default")
	if err == nil {
		t.Fatal("expected error for missing rest_url")
	}
}

func TestNewVictorOpsHandler_MissingRoutingKey(t *testing.T) {
	_, err := NewVictorOpsHandler("https://example.com", "")
	if err == nil {
		t.Fatal("expected error for missing routing_key")
	}
}

func TestNewVictorOpsHandler_Valid(t *testing.T) {
	h, err := NewVictorOpsHandler("https://example.com", "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestVictorOpsHandler_Send_Success(t *testing.T) {
	var received victorOpsPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, _ := NewVictorOpsHandler(ts.URL, "myroute")
	e := Event{
		Change: scanner.Change{Port: 8080, Type: scanner.Opened},
		Level:  LevelInfo,
		Time:   time.Now(),
	}
	if err := h.Send(e); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if received.MessageType != "INFO" {
		t.Errorf("expected INFO, got %s", received.MessageType)
	}
}

func TestVictorOpsHandler_Send_AlertSeverity(t *testing.T) {
	var received victorOpsPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received) //nolint:errcheck
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, _ := NewVictorOpsHandler(ts.URL, "myroute")
	e := Event{
		Change: scanner.Change{Port: 22, Type: scanner.Opened},
		Level:  LevelAlert,
		Time:   time.Now(),
	}
	if err := h.Send(e); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if received.MessageType != "CRITICAL" {
		t.Errorf("expected CRITICAL, got %s", received.MessageType)
	}
}

func TestVictorOpsHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h, _ := NewVictorOpsHandler(ts.URL, "myroute")
	e := Event{
		Change: scanner.Change{Port: 9090, Type: scanner.Closed},
		Level:  LevelInfo,
		Time:   time.Now(),
	}
	if err := h.Send(e); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}
