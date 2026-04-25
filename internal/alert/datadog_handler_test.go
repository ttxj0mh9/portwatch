package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewDatadogHandler_MissingKey(t *testing.T) {
	_, err := NewDatadogHandler("")
	if err == nil {
		t.Fatal("expected error for missing API key")
	}
}

func TestNewDatadogHandler_Valid(t *testing.T) {
	h, err := NewDatadogHandler("test-api-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestDatadogHandler_Send_Success(t *testing.T) {
	var received datadogEvent
	var gotAPIKey string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAPIKey = r.Header.Get("DD-API-KEY")
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	h, _ := NewDatadogHandler("my-key")
	h.url = ts.URL

	e := NewEvent(80, ChangeOpened)
	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotAPIKey != "my-key" {
		t.Errorf("expected DD-API-KEY=my-key, got %q", gotAPIKey)
	}
	if received.SourceTypeName != "portwatch" {
		t.Errorf("expected source_type_name=portwatch, got %q", received.SourceTypeName)
	}
}

func TestDatadogHandler_Send_AlertType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ev datadogEvent
		_ = json.NewDecoder(r.Body).Decode(&ev)
		if ev.AlertType != "warning" {
			t.Errorf("expected alert_type=warning for LevelAlert, got %q", ev.AlertType)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	h, _ := NewDatadogHandler("key")
	h.url = ts.URL

	e := NewEvent(22, ChangeOpened)
	e.Level = LevelAlert
	_ = h.Send(e)
}

func TestDatadogHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	h, _ := NewDatadogHandler("key")
	h.url = ts.URL

	e := NewEvent(443, ChangeOpened)
	if err := h.Send(e); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestDatadogHandler_Send_UnreachableURL(t *testing.T) {
	h, _ := NewDatadogHandler("key")
	h.url = "http://127.0.0.1:1" // nothing listening

	e := NewEvent(8080, ChangeClosed)
	if err := h.Send(e); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
