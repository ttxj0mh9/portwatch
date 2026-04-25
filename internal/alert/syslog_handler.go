package alert

import (
	"fmt"
	"log/syslog"
)

// SyslogHandler sends alerts to the system syslog daemon.
type SyslogHandler struct {
	writer *syslog.Writer
	tag    string
}

// NewSyslogHandler creates a SyslogHandler that writes to syslog with the
// given tag (e.g. "portwatch"). network and addr may be empty to use the
// local syslog socket.
func NewSyslogHandler(network, addr, tag string) (*SyslogHandler, error) {
	if tag == "" {
		tag = "portwatch"
	}
	w, err := syslog.Dial(network, addr, syslog.LOG_DAEMON|syslog.LOG_NOTICE, tag)
	if err != nil {
		return nil, fmt.Errorf("syslog: dial: %w", err)
	}
	return &SyslogHandler{writer: w, tag: tag}, nil
}

// Send writes the event to syslog. Alert-level events use LOG_WARNING;
// everything else uses LOG_NOTICE.
func (h *SyslogHandler) Send(e Event) error {
	msg := fmt.Sprintf("[%s] port %d %s", e.Level, e.Port, e.Change)
	var err error
	if e.Level == LevelAlert {
		err = h.writer.Warning(msg)
	} else {
		err = h.writer.Notice(msg)
	}
	if err != nil {
		return fmt.Errorf("syslog: write: %w", err)
	}
	return nil
}

// Close releases the underlying syslog connection.
func (h *SyslogHandler) Close() error {
	return h.writer.Close()
}
