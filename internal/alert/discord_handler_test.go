package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewDiscordHandler_MissingURL(t *testing.T) {
	_, err := NewDiscordHandler("")
	if err == nil {
		t.Fatal("expected error for empty webhook URL")
	}
}

func TestNewDiscordHandler_Valid(t *testing.T) {
	h, err := NewDiscordHandler("https://discord.com/api/webhooks/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestDiscordHandler_Send_Success(t *testing.T) {
	var received discordPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	h, _ := NewDiscordHandler(ts.URL)
	h.client = ts.Client()

	e := NewEvent("22", EventOpened)
	err := h.Send(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(received.Embeds) == 0 {
		t.Fatal("expected at least one embed in payload")
	}
	if received.Embeds[0].Color != 0x36a64f {
		t.Errorf("expected green color for info event, got %x", received.Embeds[0].Color)
	}
}

func TestDiscordHandler_Send_AlertColor(t *testing.T) {
	var received discordPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	h, _ := NewDiscordHandler(ts.URL)
	h.client = ts.Client()

	e := NewEvent("4444", EventOpened)
	err := h.Send(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(received.Embeds) == 0 {
		t.Fatal("expected embed in payload")
	}
}

func TestDiscordHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h, _ := NewDiscordHandler(ts.URL)
	h.client = ts.Client()

	e := NewEvent("80", EventOpened)
	err := h.Send(e)
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestDiscordHandler_Send_UnreachableURL(t *testing.T) {
	h, _ := NewDiscordHandler("http://127.0.0.1:1")
	e := NewEvent("80", EventClosed)
	err := h.Send(e)
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
