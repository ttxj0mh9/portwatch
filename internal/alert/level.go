package alert

import "github.com/user/portwatch/internal/scanner"

// ClassifyEvent determines the alert level for a diff event.
// If only expected ports appear in the diff it returns LevelInfo,
// otherwise it returns LevelAlert.
func ClassifyEvent(opened, closed []scanner.Port, allowed map[scanner.Port]bool) Level {
	for _, p := range opened {
		if !allowed[p] {
			return LevelAlert
		}
	}
	for _, p := range closed {
		if !allowed[p] {
			return LevelWarn
		}
	}
	if len(opened) > 0 || len(closed) > 0 {
		return LevelInfo
	}
	return LevelInfo
}

// NewEvent constructs an Event with the current timestamp and a computed level.
func NewEvent(opened, closed []scanner.Port, allowed map[scanner.Port]bool) Event {
	return Event{
		Timestamp: now(),
		Level:     ClassifyEvent(opened, closed, allowed),
		Opened:    opened,
		Closed:    closed,
	}
}
