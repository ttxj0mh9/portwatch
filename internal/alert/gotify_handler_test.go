package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewGotifyHandler_MissingURL(t *testing.T) {
	_, err := NewGotifyHandler("", "token123")
	if err == nil {
		t.Fatal("expected error for missing server URL")
	}
}

func TestNewGotifyHandler_MissingToken(t *testing.T) {
	_, err := NewGotifyHandler("http://gotify.example.com", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewGotifyHandler_Valid(t *testing.T) {
	h, err := NewGotifyHandler("http://gotify.example.com", "tok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestGotifyHandler_Send_Success(t *testing.T) {
	var received gotifyPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, _ := NewGotifyHandler(ts.URL, "testtoken")
	err := h.Send(NewEvent(8080, KindOpened))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Priority != 9 {
		t.Errorf("expected priority 9 for alert, got %d", received.Priority)
	}
}

func TestGotifyHandler_Send_InfoPriority(t *testing.T) {
	var received gotifyPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, _ := NewGotifyHandler(ts.URL, "testtoken")
	err := h.Send(NewEvent(8080, KindClosed))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Priority != 5 {
		t.Errorf("expected priority 5 for info, got %d", received.Priority)
	}
}

func TestGotifyHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	h, _ := NewGotifyHandler(ts.URL, "badtoken")
	err := h.Send(NewEvent(9090, KindOpened))
	if err == nil {
		t.Fatal("expected error on non-OK status")
	}
}

func TestGotifyHandler_Send_UnreachableURL(t *testing.T) {
	h, _ := NewGotifyHandler("http://127.0.0.1:19999", "tok")
	err := h.Send(NewEvent(80, KindOpened))
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
