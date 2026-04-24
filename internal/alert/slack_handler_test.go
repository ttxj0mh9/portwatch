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

func TestSlackHandler_Send_Success(t *testing.T) {
	var received string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var p struct {
			Text string `json:"text"`
		}
		_ = json.Unmarshal(body, &p)
		received = p.Text
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	h := NewSlackHandler(srv.URL)
	e := NewEvent(9000, EventOpened, time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))

	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(received, "9000") {
		t.Errorf("expected port 9000 in message, got: %s", received)
	}
}

func TestSlackHandler_Send_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	h := NewSlackHandler(srv.URL)
	e := NewEvent(9001, EventClosed, time.Now())

	if err := h.Send(e); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestSlackHandler_Send_UnreachableURL(t *testing.T) {
	h := NewSlackHandler("http://127.0.0.1:1")
	e := NewEvent(9002, EventOpened, time.Now())

	if err := h.Send(e); err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}
