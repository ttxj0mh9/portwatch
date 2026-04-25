package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewTelegramHandler_MissingToken(t *testing.T) {
	_, err := NewTelegramHandler("", "123456")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewTelegramHandler_MissingChatID(t *testing.T) {
	_, err := NewTelegramHandler("bot-token", "")
	if err == nil {
		t.Fatal("expected error for missing chat ID")
	}
}

func TestNewTelegramHandler_Valid(t *testing.T) {
	h, err := NewTelegramHandler("bot-token", "123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestTelegramHandler_Send_Success(t *testing.T) {
	var received map[string]string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, _ := NewTelegramHandler("test-token", "chat-99")
	h.client = &http.Client{
		Transport: rewriteTransport(telegramAPIBase+"test-token", ts.URL),
	}

	e := NewEvent(8080, EventOpened, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["chat_id"] != "chat-99" {
		t.Errorf("expected chat_id 'chat-99', got %q", received["chat_id"])
	}
	if received["text"] == "" {
		t.Error("expected non-empty text")
	}
	if received["parse_mode"] != "Markdown" {
		t.Errorf("expected parse_mode 'Markdown', got %q", received["parse_mode"])
	}
}

func TestTelegramHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	h, _ := NewTelegramHandler("bad-token", "chat-99")
	h.client = &http.Client{
		Transport: rewriteTransport(telegramAPIBase+"bad-token", ts.URL),
	}

	e := NewEvent(9090, EventClosed, time.Now())
	if err := h.Send(e); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestTelegramHandler_Send_UnreachableURL(t *testing.T) {
	h, _ := NewTelegramHandler("tok", "cid")
	h.client = &http.Client{}
	// point at a URL that will refuse connections
	h.token = "tok"

	e := NewEvent(443, EventOpened, time.Now())
	// Override client to use a bad transport
	h.client = &http.Client{Transport: rewriteTransport(telegramAPIBase+"tok", "http://127.0.0.1:1")}
	if err := h.Send(e); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
