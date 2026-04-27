package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewMatrixHandler_MissingHomeserver(t *testing.T) {
	_, err := NewMatrixHandler("", "token", "!room:matrix.org")
	if err == nil || !strings.Contains(err.Error(), "homeserver") {
		t.Fatalf("expected homeserver error, got %v", err)
	}
}

func TestNewMatrixHandler_MissingToken(t *testing.T) {
	_, err := NewMatrixHandler("https://matrix.org", "", "!room:matrix.org")
	if err == nil || !strings.Contains(err.Error(), "access token") {
		t.Fatalf("expected access token error, got %v", err)
	}
}

func TestNewMatrixHandler_MissingRoomID(t *testing.T) {
	_, err := NewMatrixHandler("https://matrix.org", "token", "")
	if err == nil || !strings.Contains(err.Error(), "room ID") {
		t.Fatalf("expected room ID error, got %v", err)
	}
}

func TestNewMatrixHandler_Valid(t *testing.T) {
	h, err := NewMatrixHandler("https://matrix.org", "token", "!room:matrix.org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestMatrixHandler_Send_Success(t *testing.T) {
	var received map[string]string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			t.Errorf("expected Bearer token, got %q", auth)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"event_id":"$abc"}`))
	}))
	defer srv.Close()

	h, _ := NewMatrixHandler(srv.URL, "mytoken", "!testroom:example.com")
	e := NewEvent(9200, StateOpened)

	if err := h.Send(e); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if received["msgtype"] != "m.text" {
		t.Errorf("expected msgtype m.text, got %q", received["msgtype"])
	}
	if received["body"] == "" {
		t.Error("expected non-empty body")
	}
}

func TestMatrixHandler_Send_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	h, _ := NewMatrixHandler(srv.URL, "badtoken", "!room:example.com")
	e := NewEvent(443, StateClosed)

	if err := h.Send(e); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestMatrixHandler_Send_UnreachableURL(t *testing.T) {
	h, _ := NewMatrixHandler("http://127.0.0.1:1", "token", "!room:example.com")
	e := NewEvent(80, StateOpened)
	if err := h.Send(e); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
