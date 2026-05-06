package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func validZendeskConfig(url string) config.ZendeskConfig {
	return config.ZendeskConfig{
		Subdomain: "testcorp",
		Email:     "admin@example.com",
		APIToken:  "secret-token",
	}
}

func TestNewZendeskHandler_MissingSubdomain(t *testing.T) {
	_, err := NewZendeskHandler(config.ZendeskConfig{
		Email:    "a@b.com",
		APIToken: "tok",
	})
	if err == nil {
		t.Fatal("expected error for missing subdomain")
	}
}

func TestNewZendeskHandler_MissingEmail(t *testing.T) {
	_, err := NewZendeskHandler(config.ZendeskConfig{
		Subdomain: "corp",
		APIToken:  "tok",
	})
	if err == nil {
		t.Fatal("expected error for missing email")
	}
}

func TestNewZendeskHandler_MissingAPIToken(t *testing.T) {
	_, err := NewZendeskHandler(config.ZendeskConfig{
		Subdomain: "corp",
		Email:     "a@b.com",
	})
	if err == nil {
		t.Fatal("expected error for missing api_token")
	}
}

func TestNewZendeskHandler_Valid(t *testing.T) {
	h, err := NewZendeskHandler(validZendeskConfig(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestZendeskHandler_Send_Success(t *testing.T) {
	var received zendeskTicket
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	h, _ := NewZendeskHandler(validZendeskConfig(ts.URL))
	h.client = rewriteTransport(h.client, "testcorp.zendesk.com", ts.URL)

	event := NewEvent(8080, EventOpened)
	if err := h.Send(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Ticket.Priority != "normal" {
		t.Errorf("expected priority normal, got %s", received.Ticket.Priority)
	}
}

func TestZendeskHandler_Send_AlertPriority(t *testing.T) {
	var received zendeskTicket
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	h, _ := NewZendeskHandler(validZendeskConfig(ts.URL))
	h.client = rewriteTransport(h.client, "testcorp.zendesk.com", ts.URL)

	event := NewEvent(4444, EventOpened)
	event.Level = LevelAlert
	if err := h.Send(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Ticket.Priority != "high" {
		t.Errorf("expected priority high, got %s", received.Ticket.Priority)
	}
}

func TestZendeskHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h, _ := NewZendeskHandler(validZendeskConfig(ts.URL))
	h.client = rewriteTransport(h.client, "testcorp.zendesk.com", ts.URL)

	event := NewEvent(8080, EventOpened)
	if err := h.Send(event); err == nil {
		t.Fatal("expected error for non-201 status")
	}
}
