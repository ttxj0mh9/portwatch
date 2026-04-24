package monitor

import (
	"fmt"
	"strings"
)

// FormatAlert returns a human-readable string describing an Alert.
func FormatAlert(a Alert) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("[%s] Port change detected\n", a.Timestamp.Format("2006-01-02 15:04:05")))

	if len(a.Opened) > 0 {
		parts := make([]string, len(a.Opened))
		for i, p := range a.Opened {
			parts[i] = p.String()
		}
		sb.WriteString(fmt.Sprintf("  OPENED: %s\n", strings.Join(parts, ", ")))
	}

	if len(a.Closed) > 0 {
		parts := make([]string, len(a.Closed))
		for i, p := range a.Closed {
			parts[i] = p.String()
		}
		sb.WriteString(fmt.Sprintf("  CLOSED: %s\n", strings.Join(parts, ", ")))
	}

	return sb.String()
}

// Summary returns a compact one-line summary of an Alert.
func Summary(a Alert) string {
	return fmt.Sprintf("%d opened, %d closed at %s",
		len(a.Opened),
		len(a.Closed),
		a.Timestamp.Format("15:04:05"),
	)
}
