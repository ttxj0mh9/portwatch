package alert

import (
	"fmt"
	"net/smtp"
	"strings"
)

// EmailConfig holds SMTP configuration for sending email alerts.
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
}

// emailHandler sends alerts via SMTP email.
type emailHandler struct {
	cfg  EmailConfig
	auth smtp.Auth
}

// NewEmailHandler creates a new email alert handler.
func NewEmailHandler(cfg EmailConfig) (Handler, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("email handler: SMTP host is required")
	}
	if len(cfg.To) == 0 {
		return nil, fmt.Errorf("email handler: at least one recipient is required")
	}
	var auth smtp.Auth
	if cfg.Username != "" {
		auth = smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	}
	return &emailHandler{cfg: cfg, auth: auth}, nil
}

// Send delivers an alert event as an email message.
func (h *emailHandler) Send(e Event) error {
	addr := fmt.Sprintf("%s:%d", h.cfg.Host, h.cfg.Port)
	subject := fmt.Sprintf("[portwatch] %s: port %d %s", e.Level, e.Port, e.Kind)
	body := fmt.Sprintf(
		"To: %s\r\nFrom: %s\r\nSubject: %s\r\n\r\n%s\r\n",
		strings.Join(h.cfg.To, ", "),
		h.cfg.From,
		subject,
		formatEmailBody(e),
	)
	return smtp.SendMail(addr, h.auth, h.cfg.From, h.cfg.To, []byte(body))
}

func formatEmailBody(e Event) string {
	return fmt.Sprintf(
		"Port:      %d\nEvent:     %s\nLevel:     %s\nTimestamp: %s\n",
		e.Port,
		e.Kind,
		e.Level,
		e.Time.Format(timeFormat),
	)
}
