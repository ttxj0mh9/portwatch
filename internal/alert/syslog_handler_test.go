package alert

import (
	"testing"
	"time"
)

// fakeSyslogWriter captures messages written to it so tests do not need a
// real syslog socket.
type fakeSyslogWriter struct {
	notices  []string
	warnings []string
}

func (f *fakeSyslogWriter) Notice(m string) error  { f.notices = append(f.notices, m); return nil }
func (f *fakeSyslogWriter) Warning(m string) error { f.warnings = append(f.warnings, m); return nil }
func (f *fakeSyslogWriter) Close() error           { return nil }

// syslogWriterIface is the subset of *syslog.Writer used by SyslogHandler,
// extracted so we can inject a fake in tests.
type syslogWriterIface interface {
	Notice(string) error
	Warning(string) error
	Close() error
}

// newSyslogHandlerFromWriter creates a SyslogHandler backed by an arbitrary
// syslogWriterIface — used only in tests.
func newSyslogHandlerFromWriter(w syslogWriterIface) *SyslogHandler {
	// We store the concrete *syslog.Writer in the struct; for tests we use
	// a small adapter to bridge the interface.
	return &SyslogHandler{writer: nil, tag: "test", iface: w}
}

func TestSyslogHandler_Send_Info(t *testing.T) {
	fw := &fakeSyslogWriter{}
	h := newSyslogHandlerFromWriter(fw)

	e := Event{Port: 8080, Change: "opened", Level: LevelInfo, Time: time.Now()}
	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fw.notices) != 1 {
		t.Fatalf("expected 1 notice, got %d", len(fw.notices))
	}
	if len(fw.warnings) != 0 {
		t.Fatalf("expected 0 warnings, got %d", len(fw.warnings))
	}
}

func TestSyslogHandler_Send_Alert(t *testing.T) {
	fw := &fakeSyslogWriter{}
	h := newSyslogHandlerFromWriter(fw)

	e := Event{Port: 22, Change: "closed", Level: LevelAlert, Time: time.Now()}
	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fw.warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(fw.warnings))
	}
	if len(fw.notices) != 0 {
		t.Fatalf("expected 0 notices, got %d", len(fw.notices))
	}
}

func TestNewSyslogHandler_InvalidAddr(t *testing.T) {
	_, err := NewSyslogHandler("tcp", "127.0.0.1:0", "portwatch")
	if err == nil {
		t.Fatal("expected error for unreachable syslog address")
	}
}
