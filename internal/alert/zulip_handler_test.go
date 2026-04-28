package alert

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewZulipHandler_MissingBaseURL(t *testing.T) {
	_, err := NewZulipHandler("", "bot@example.com", "key", "general", "alerts")
	if err == nil {
		t.Fatal("expected error for missing baseURL")
	}
}

func TestNewZulipHandler_MissingEmail(t *testing.T) {
	_, err := NewZulipHandler("https://org.zulipchat.com", "", "key", "general", "alerts")
	if err == nil {
		t.Fatal("expected error for missing email")
	}
}

func TestNewZulipHandler_MissingAPIKey(t *testing.T) {
	_, err := NewZulipHandler("https://org.zulipchat.com", "bot@example.com", "", "general", "alerts")
	if err == nil {
		t.Fatal("expected error for missing API key")
	}
}

func TestNewZulipHandler_MissingStream(t *testing.T) {
	_, err := NewZulipHandler("https://org.zulipchat.com", "bot@example.com", "key", "", "alerts")
	if err == nil {
		t.Fatal("expected error for missing stream")
	}
}

func TestNewZulipHandler_DefaultTopic(t *testing.T) {
	h, err := NewZulipHandler("https://org.zulipchat.com", "bot@example.com", "key", "general", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.topic != "portwatch alerts" {
		t.Errorf("expected default topic, got %q", h.topic)
	}
}

func TestZulipHandler_Send_Success(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	h, err := NewZulipHandler(server.URL, "bot@example.com", "secret", "general", "alerts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e := Event{
		Port:      8080,
		State:     "opened",
		Level:     LevelAlert,
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	if err := h.Send(e); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if gotAuth == "" {
		t.Error("expected Authorization header to be set")
	}
}

func TestZulipHandler_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	h, err := NewZulipHandler(server.URL, "bot@example.com", "bad-key", "general", "alerts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e := Event{Port: 443, State: "closed", Level: LevelInfo}
	if err := h.Send(e); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestZulipHandler_Send_UnreachableURL(t *testing.T) {
	h, err := NewZulipHandler("http://127.0.0.1:19999", "bot@example.com", "key", "general", "alerts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := Event{Port: 22, State: "opened", Level: LevelAlert}
	if err := h.Send(e); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
