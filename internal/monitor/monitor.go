package monitor

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Alert represents a change detected in open ports.
type Alert struct {
	Opened []scanner.Port
	Closed []scanner.Port
	Timestamp time.Time
}

// Monitor periodically scans ports and emits alerts on changes.
type Monitor struct {
	scanner  *scanner.TCPScanner
	ports    []int
	interval time.Duration
	Alerts   chan Alert
}

// New creates a new Monitor.
func New(s *scanner.TCPScanner, ports []int, interval time.Duration) *Monitor {
	return &Monitor{
		scanner:  s,
		ports:    ports,
		interval: interval,
		Alerts:   make(chan Alert, 16),
	}
}

// Run starts the monitoring loop. It blocks until ctx is cancelled.
func (m *Monitor) Run(ctx context.Context) {
	prev, err := m.scanner.Scan(m.ports)
	if err != nil {
		log.Printf("portwatch: initial scan error: %v", err)
	}

	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()
	defer close(m.Alerts)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			current, err := m.scanner.Scan(m.ports)
			if err != nil {
				log.Printf("portwatch: scan error: %v", err)
				continue
			}

			diff := scanner.Diff(prev, current)
			if len(diff.Opened) > 0 || len(diff.Closed) > 0 {
				m.sendAlert(Alert{
					Opened:    diff.Opened,
					Closed:    diff.Closed,
					Timestamp: time.Now(),
				})
			}
			prev = current
		}
	}
}

// sendAlert attempts to send an alert on the Alerts channel. If the channel
// buffer is full, the alert is dropped and a warning is logged to avoid
// blocking the monitoring loop.
func (m *Monitor) sendAlert(a Alert) {
	select {
	case m.Alerts <- a:
	default:
		log.Printf("portwatch: alert channel full, dropping alert (opened=%d, closed=%d)",
			len(a.Opened), len(a.Closed))
	}
}
