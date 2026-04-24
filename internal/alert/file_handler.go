package alert

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// FileHandler writes alert events to a log file, rotating by date.
type FileHandler struct {
	mu      sync.Mutex
	path    string
	current string
	file    *os.File
}

// NewFileHandler creates a FileHandler that writes to the given file path.
func NewFileHandler(path string) *FileHandler {
	return &FileHandler{path: path}
}

// Send appends the event to the log file.
func (f *FileHandler) Send(e Event) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if err := f.ensureOpen(); err != nil {
		return err
	}

	for _, p := range e.Opened {
		if _, err := fmt.Fprintf(f.file, "%s [%s] port opened: %s\n",
			e.Timestamp.Format(time.RFC3339), e.Level, p); err != nil {
			return err
		}
	}
	for _, p := range e.Closed {
		if _, err := fmt.Fprintf(f.file, "%s [%s] port closed: %s\n",
			e.Timestamp.Format(time.RFC3339), e.Level, p); err != nil {
			return err
		}
	}
	return nil
}

// Close releases the underlying file handle.
func (f *FileHandler) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.file != nil {
		err := f.file.Close()
		f.file = nil
		return err
	}
	return nil
}

func (f *FileHandler) ensureOpen() error {
	if f.file == nil {
		var err error
		f.file, err = os.OpenFile(f.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return fmt.Errorf("alert file handler: %w", err)
		}
	}
	return nil
}
