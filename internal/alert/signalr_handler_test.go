package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewSignalRHandler_MissingURL(t *testing.T) {
	_, err := NewSignalRHandler("", "key")
	if err == nil {
		t.Fatal("expected error for missing hub URL")
	}
}

func TestNewSignalRHandler_MissingKey(t *testing.T) {
	_, err := NewSignalRHandler("http://example.com", "")
	if err == nil {
		t.Fatal("expected error for missing access key")
	}
}

func TestNewSignalRHandler_Valid(t *testing.T) {
	h, err := NewSignalRHandler("http://example.com", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestSignalRHandler_Send_Success(t *testing.T) {
	var captured signalRPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer testkey" {
			t.Errorf("missing or wrong Authorization header")
		}
		if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, _ := NewSignalRHandler(ts.URL, "testkey")
	e := Event{Port: 8080, Kind: Opened, Level: LevelAlert, Time: time.Now()}
	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured.Target != "portwatch" {
		t.Errorf("expected target 'portwatch', got %q", captured.Target)
	}
	if len(captured.Arguments) != 1 || captured.Arguments[0].Port != 8080 {
		t.Errorf("unexpected arguments: %+v", captured.Arguments)
	}
}

func TestSignalRHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	h, _ := NewSignalRHandler(ts.URL, "badkey")
	e := Event{Port: 443, Kind: Closed, Level: LevelInfo, Time: time.Now()}
	if err := h.Send(e); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestSignalRHandler_Send_UnreachableURL(t *testing.T) {
	h, _ := NewSignalRHandler("http://127.0.0.1:1", "key")
	e := Event{Port: 22, Kind: Opened, Level: LevelAlert, Time: time.Now()}
	if err := h.Send(e); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
