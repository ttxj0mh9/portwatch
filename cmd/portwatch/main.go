// main is the entry point for the portwatch CLI daemon.
// It loads configuration, wires up the scanner, monitor, alert dispatcher,
// and snapshot store, then runs the monitoring loop until interrupted.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/portwatch/internal/alert"
	"github.com/yourorg/portwatch/internal/config"
	"github.com/yourorg/portwatch/internal/monitor"
	"github.com/yourorg/portwatch/internal/scanner"
	"github.com/yourorg/portwatch/internal/snapshot"
)

func main() {
	cfgPath := flag.String("config", "", "path to config file (default: built-in defaults)")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Build the alert dispatcher from config.
	dispatcher, err := buildDispatcher(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: failed to build alert dispatcher: %v\n", err)
		os.Exit(1)
	}

	// Set up the rotating snapshot store so we keep a history of port states.
	store, err := snapshot.NewRotatingStore(cfg.SnapshotPath, cfg.SnapshotKeep)
	if err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: failed to init snapshot store: %v\n", err)
		os.Exit(1)
	}

	scn := scanner.NewTCPScanner(cfg.Ports, cfg.ScanTimeout)
	mon := monitor.New(scn, store, dispatcher, cfg.Interval)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Printf("portwatch: starting — watching %d port(s) every %s", len(cfg.Ports), cfg.Interval)

	if err := mon.Run(ctx); err != nil {
		log.Printf("portwatch: monitor exited: %v", err)
	}

	log.Println("portwatch: stopped")
}

// buildDispatcher constructs an alert.Dispatcher wired with all handlers
// that are enabled in the configuration.
func buildDispatcher(cfg *config.Config) (*alert.Dispatcher, error) {
	var handlers []alert.Handler

	// Log handler is always active.
	handlers = append(handlers, alert.NewLogHandler(log.Default()))

	// Optional file handler.
	if cfg.Alert.File.Path != "" {
		fh, err := alert.NewFileHandler(cfg.Alert.File.Path)
		if err != nil {
			return nil, fmt.Errorf("file handler: %w", err)
		}
		handlers = append(handlers, fh)
	}

	// Optional webhook handler.
	if cfg.Alert.Webhook.URL != "" {
		handlers = append(handlers, alert.NewWebhookHandler(cfg.Alert.Webhook.URL))
	}

	// Optional Slack handler.
	if cfg.Alert.Slack.WebhookURL != "" {
		handlers = append(handlers, alert.NewSlackHandler(cfg.Alert.Slack.WebhookURL))
	}

	// Optional PagerDuty handler.
	if cfg.Alert.PagerDuty.IntegrationKey != "" {
		pd, err := alert.NewPagerDutyHandler(cfg.Alert.PagerDuty.IntegrationKey)
		if err != nil {
			return nil, fmt.Errorf("pagerduty handler: %w", err)
		}
		handlers = append(handlers, pd)
	}

	// Optional email handler.
	if cfg.Alert.Email.Host != "" {
		eh, err := alert.NewEmailHandler(
			cfg.Alert.Email.Host,
			cfg.Alert.Email.Port,
			cfg.Alert.Email.From,
			cfg.Alert.Email.To,
			cfg.Alert.Email.Username,
			cfg.Alert.Email.Password,
		)
		if err != nil {
			return nil, fmt.Errorf("email handler: %w", err)
		}
		handlers = append(handlers, eh)
	}

	return alert.NewDispatcher(handlers...), nil
}
