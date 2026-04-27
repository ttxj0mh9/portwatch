package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func TestNewRocketChatHandler_MissingURL(t *testing.T) {
	_, err := NewRocketChatHandler("", "#alerts", "portwatch")
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestNewRocketChatHandler_Valid(t *testing.T) {
	h, err := NewRocketChatHandler("http://example.com/hook", "#alerts", "portwatch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestRocketChatHandler_Send_Success(t *testing.T) {
	var captured map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, _ := NewRocketChatHandler(ts.URL, "#alerts", "portwatch")
	snap := snapshot.New([]uint16{8080})
	event := NewEvent(snap, EventOpened)

	if err := h.Send(event, snap); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured["text"] == nil {
		t.Error("expected text field in payload")
	}
	if captured["channel"] != "#alerts" {
		t.Errorf("expected channel #alerts, got %v", captured["channel"])
	}
}

func TestRocketChatHandler_Send_AlertEmoji(t *testing.T) {
	var captured map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&captured)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, _ := NewRocketChatHandler(ts.URL, "", "")
	snap := snapshot.New([]uint16{22})
	event := NewEvent(snap, EventOpened)
	event.Level = LevelAlert

	h.Send(event, snap)

	text, _ := captured["text"].(string)
	if len(text) == 0 {
		t.Error("expected non-empty text")
	}
}

func TestRocketChatHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h, _ := NewRocketChatHandler(ts.URL, "#alerts", "portwatch")
	snap := snapshot.New([]uint16{9090})
	event := NewEvent(snap, EventClosed)

	if err := h.Send(event, snap); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestRocketChatHandler_Send_UnreachableURL(t *testing.T) {
	h, _ := NewRocketChatHandler("http://127.0.0.1:1", "", "")
	snap := snapshot.New([]uint16{3000})
	event := NewEvent(snap, EventOpened)

	if err := h.Send(event, snap); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
