package alert

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestNoopHandler_Send_ReturnsNil(t *testing.T) {
	h := NewNoopHandler()
	e := NewEvent(8080, "opened")
	if err := h.Send(e); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestNoopHandler_Send_Silent(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	h := NewNoopHandler(WithNoopLogger(logger))
	_ = h.Send(NewEvent(443, "closed"))

	if buf.Len() != 0 {
		t.Fatalf("expected no log output in silent mode, got: %q", buf.String())
	}
}

func TestNoopHandler_Send_Verbose(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	h := NewNoopHandler(WithNoopLogger(logger), WithNoopVerbose(true))
	e := NewEvent(22, "opened")
	_ = h.Send(e)

	out := buf.String()
	if !strings.Contains(out, "noop") {
		t.Errorf("expected log line to contain 'noop', got: %q", out)
	}
	if !strings.Contains(out, "22") {
		t.Errorf("expected log line to contain port '22', got: %q", out)
	}
}

func TestNoopHandler_Send_MultipleEvents(t *testing.T) {
	h := NewNoopHandler()
	events := []Event{
		NewEvent(80, "opened"),
		NewEvent(443, "closed"),
		NewEvent(8080, "opened"),
	}
	for _, e := range events {
		if err := h.Send(e); err != nil {
			t.Errorf("unexpected error for port %d: %v", e.Port, err)
		}
	}
}
