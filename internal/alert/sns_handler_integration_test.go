//go:build integration
// +build integration

package alert

import (
	"os"
	"testing"
	"time"
)

// TestSNSHandler_RealPublish requires AWS credentials and a real SNS topic ARN
// set via the PORTWATCH_TEST_SNS_ARN environment variable.
// Run with: go test -tags integration ./internal/alert/...
func TestSNSHandler_RealPublish(t *testing.T) {
	arn := os.Getenv("PORTWATCH_TEST_SNS_ARN")
	if arn == "" {
		t.Skip("PORTWATCH_TEST_SNS_ARN not set; skipping integration test")
	}

	h, err := NewSNSHandler(arn)
	if err != nil {
		t.Fatalf("failed to create SNS handler: %v", err)
	}

	e := Event{
		Kind:      KindOpened,
		Port:      9999,
		Timestamp: time.Now(),
	}

	if err := h.Send(e); err != nil {
		t.Fatalf("Send failed: %v", err)
	}
}
