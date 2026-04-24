package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event holds the details of a single alert event.
type Event struct {
	Timestamp time.Time
	Level     Level
	Opened    []scanner.Port
	Closed    []scanner.Port
}

// Handler is the interface for alert output destinations.
type Handler interface {
	Send(e Event) error
}

// Dispatcher fans out an event to multiple handlers.
type Dispatcher struct {
	handlers []Handler
}

// NewDispatcher creates a Dispatcher with the given handlers.
func NewDispatcher(handlers ...Handler) *Dispatcher {
	return &Dispatcher{handlers: handlers}
}

// Dispatch sends the event to all registered handlers.
func (d *Dispatcher) Dispatch(e Event) error {
	var lastErr error
	for _, h := range d.handlers {
		if err := h.Send(e); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// LogHandler writes alert events to an io.Writer.
type LogHandler struct {
	w io.Writer
}

// NewLogHandler creates a LogHandler writing to w.
// Pass nil to default to os.Stdout.
func NewLogHandler(w io.Writer) *LogHandler {
	if w == nil {
		w = os.Stdout
	}
	return &LogHandler{w: w}
}

// Send formats and writes the event to the underlying writer.
func (l *LogHandler) Send(e Event) error {
	for _, p := range e.Opened {
		_, err := fmt.Fprintf(l.w, "%s [%s] port opened: %s\n",
			e.Timestamp.Format(time.RFC3339), e.Level, p)
		if err != nil {
			return err
		}
	}
	for _, p := range e.Closed {
		_, err := fmt.Fprintf(l.w, "%s [%s] port closed: %s\n",
			e.Timestamp.Format(time.RFC3339), e.Level, p)
		if err != nil {
			return err
		}
	}
	return nil
}
