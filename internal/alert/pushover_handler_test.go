package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewPushoverHandler_MissingUserKey(t *testing.T) {
	_, err := NewPushoverHandler("", "token123")
	if err == nil {
		t.Fatal("expected error for missing user key")
	}
}

func TestNewPushoverHandler_MissingAPIToken(t *testing.T) {
	_, err := NewPushoverHandler("userkey123", "")
	if err == nil {
		t.Fatal("expected error for missing api token")
	}
}

func TestNewPushoverHandler_Valid(t *testing.T) {
	h, err := NewPushoverHandler("userkey123", "token123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestPushoverHandler_Send_Success(t *testing.T) {
	var received map[string][]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		received = map[string][]string(r.Form)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": 1})
	}))
	defer server.Close()

	h, _ := NewPushoverHandler("ukey", "atoken")
	h.client = server.Client()
	h.client = &http.Client{Transport: rewriteTransport(server.URL)}

	event := Event{Port: 8080, Change: "opened", Level: LevelInfo, Time: time.Now()}
	if err := h.Send(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["token"][0] != "atoken" {
		t.Errorf("expected token atoken, got %s", received["token"][0])
	}
	if received["user"][0] != "ukey" {
		t.Errorf("expected user ukey, got %s", received["user"][0])
	}
}

func TestPushoverHandler_Send_AlertPriority(t *testing.T) {
	var received map[string][]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		received = map[string][]string(r.Form)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": 1})
	}))
	defer server.Close()

	h, _ := NewPushoverHandler("ukey", "atoken")
	h.client = &http.Client{Transport: rewriteTransport(server.URL)}

	event := Event{Port: 22, Change: "opened", Level: LevelAlert, Time: time.Now()}
	_ = h.Send(event)

	if received["priority"][0] != "1" {
		t.Errorf("expected priority 1 for alert, got %s", received["priority"][0])
	}
}

func TestPushoverHandler_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	h, _ := NewPushoverHandler("ukey", "atoken")
	h.client = &http.Client{Transport: rewriteTransport(server.URL)}

	event := Event{Port: 9090, Change: "closed", Level: LevelInfo, Time: time.Now()}
	if err := h.Send(event); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}
