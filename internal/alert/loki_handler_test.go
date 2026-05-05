package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewLokiHandler_MissingURL(t *testing.T) {
	_, err := NewLokiHandler("", nil)
	if err == nil {
		t.Fatal("expected error for missing push URL")
	}
}

func TestNewLokiHandler_DefaultLabels(t *testing.T) {
	h, err := NewLokiHandler("http://localhost:3100/loki/api/v1/push", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.labels["app"] != "portwatch" {
		t.Errorf("expected default label app=portwatch, got %v", h.labels)
	}
}

func TestNewLokiHandler_Valid(t *testing.T) {
	h, err := NewLokiHandler("http://loki:3100/loki/api/v1/push", map[string]string{"env": "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.pushURL == "" {
		t.Error("expected pushURL to be set")
	}
}

func TestLokiHandler_Send_Success(t *testing.T) {
	var received []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	h, _ := NewLokiHandler(ts.URL, map[string]string{"app": "portwatch"})
	ev := NewEvent(8080, true, time.Unix(0, 1_000_000_000))

	if err := h.Send(ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload lokiPushPayload
	if err := json.Unmarshal(received, &payload); err != nil {
		t.Fatalf("could not unmarshal payload: %v", err)
	}
	if len(payload.Streams) != 1 {
		t.Fatalf("expected 1 stream, got %d", len(payload.Streams))
	}
	if payload.Streams[0].Stream["app"] != "portwatch" {
		t.Errorf("expected label app=portwatch")
	}
	if len(payload.Streams[0].Values) != 1 {
		t.Fatalf("expected 1 log entry")
	}
	if !strings.Contains(payload.Streams[0].Values[0][1], "8080") {
		t.Errorf("expected log line to mention port 8080")
	}
}

func TestLokiHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h, _ := NewLokiHandler(ts.URL, nil)
	ev := NewEvent(9090, false, time.Now())
	if err := h.Send(ev); err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestLokiHandler_Send_UnreachableURL(t *testing.T) {
	h, _ := NewLokiHandler("http://127.0.0.1:19999/loki/api/v1/push", nil)
	ev := NewEvent(443, true, time.Now())
	if err := h.Send(ev); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
