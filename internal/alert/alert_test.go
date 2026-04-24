package alert

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func fixedTime() time.Time {
	return time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
}

func TestLogHandler_Send_Opened(t *testing.T) {
	now = fixedTime
	var buf bytes.Buffer
	h := NewLogHandler(&buf)
	e := Event{
		Timestamp: fixedTime(),
		Level:     LevelAlert,
		Opened:    []scanner.Port{8080},
	}
	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "port opened: 8080/tcp") {
		t.Errorf("expected 'port opened: 8080/tcp' in output, got: %q", got)
	}
	if !strings.Contains(got, "ALERT") {
		t.Errorf("expected level ALERT in output, got: %q", got)
	}
}

func TestLogHandler_Send_Closed(t *testing.T) {
	now = fixedTime
	var buf bytes.Buffer
	h := NewLogHandler(&buf)
	e := Event{
		Timestamp: fixedTime(),
		Level:     LevelWarn,
		Closed:    []scanner.Port{22},
	}
	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "port closed: 22/tcp") {
		t.Errorf("expected 'port closed: 22/tcp' in output, got: %q", got)
	}
}

func TestClassifyEvent_Alert(t *testing.T) {
	allowed := map[scanner.Port]bool{80: true}
	level := ClassifyEvent([]scanner.Port{8080}, nil, allowed)
	if level != LevelAlert {
		t.Errorf("expected ALERT, got %s", level)
	}
}

func TestClassifyEvent_Info(t *testing.T) {
	allowed := map[scanner.Port]bool{8080: true}
	level := ClassifyEvent([]scanner.Port{8080}, nil, allowed)
	if level != LevelInfo {
		t.Errorf("expected INFO, got %s", level)
	}
}

func TestClassifyEvent_Warn(t *testing.T) {
	allowed := map[scanner.Port]bool{22: true}
	level := ClassifyEvent(nil, []scanner.Port{22}, allowed)
	if level != LevelWarn {
		t.Errorf("expected WARN, got %s", level)
	}
}

func TestDispatcher_Dispatch(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	d := NewDispatcher(NewLogHandler(&buf1), NewLogHandler(&buf2))
	e := Event{
		Timestamp: fixedTime(),
		Level:     LevelInfo,
		Opened:    []scanner.Port{3000},
	}
	if err := d.Dispatch(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf1.String(), "3000/tcp") {
		t.Error("handler 1 did not receive event")
	}
	if !strings.Contains(buf2.String(), "3000/tcp") {
		t.Error("handler 2 did not receive event")
	}
}
