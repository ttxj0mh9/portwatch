package scanner_test

import (
	"testing"

	"github.com/yourusername/portwatch/internal/scanner"
)

func TestDiff_NewPortOpened(t *testing.T) {
	prev := []scanner.Port{}
	curr := []scanner.Port{
		{Number: 8080, Protocol: "tcp"},
	}

	result := scanner.Diff(prev, curr)

	if len(result.Opened) != 1 {
		t.Fatalf("expected 1 opened port, got %d", len(result.Opened))
	}
	if result.Opened[0].Number != 8080 {
		t.Errorf("expected opened port 8080, got %d", result.Opened[0].Number)
	}
	if len(result.Closed) != 0 {
		t.Errorf("expected 0 closed ports, got %d", len(result.Closed))
	}
}

func TestDiff_PortClosed(t *testing.T) {
	prev := []scanner.Port{
		{Number: 22, Protocol: "tcp"},
		{Number: 80, Protocol: "tcp"},
	}
	curr := []scanner.Port{
		{Number: 80, Protocol: "tcp"},
	}

	result := scanner.Diff(prev, curr)

	if len(result.Closed) != 1 {
		t.Fatalf("expected 1 closed port, got %d", len(result.Closed))
	}
	if result.Closed[0].Number != 22 {
		t.Errorf("expected closed port 22, got %d", result.Closed[0].Number)
	}
	if len(result.Opened) != 0 {
		t.Errorf("expected 0 opened ports, got %d", len(result.Opened))
	}
}

func TestDiff_NoChange(t *testing.T) {
	ports := []scanner.Port{
		{Number: 443, Protocol: "tcp"},
		{Number: 8443, Protocol: "tcp"},
	}

	result := scanner.Diff(ports, ports)

	if len(result.Opened) != 0 {
		t.Errorf("expected 0 opened ports, got %d", len(result.Opened))
	}
	if len(result.Closed) != 0 {
		t.Errorf("expected 0 closed ports, got %d", len(result.Closed))
	}
}

func TestDiff_MultipleChanges(t *testing.T) {
	prev := []scanner.Port{
		{Number: 22, Protocol: "tcp"},
		{Number: 80, Protocol: "tcp"},
		{Number: 3306, Protocol: "tcp"},
	}
	curr := []scanner.Port{
		{Number: 22, Protocol: "tcp"},
		{Number: 443, Protocol: "tcp"},
		{Number: 8080, Protocol: "tcp"},
	}

	result := scanner.Diff(prev, curr)

	if len(result.Opened) != 2 {
		t.Fatalf("expected 2 opened ports, got %d", len(result.Opened))
	}
	if len(result.Closed) != 2 {
		t.Fatalf("expected 2 closed ports, got %d", len(result.Closed))
	}

	// Verify specific ports
	openedSet := make(map[int]bool)
	for _, p := range result.Opened {
		openedSet[p.Number] = true
	}
	if !openedSet[443] || !openedSet[8080] {
		t.Errorf("expected opened ports 443 and 8080, got %v", result.Opened)
	}

	closedSet := make(map[int]bool)
	for _, p := range result.Closed {
		closedSet[p.Number] = true
	}
	if !closedSet[80] || !closedSet[3306] {
		t.Errorf("expected closed ports 80 and 3306, got %v", result.Closed)
	}
}

func TestDiff_BothEmpty(t *testing.T) {
	result := scanner.Diff([]scanner.Port{}, []scanner.Port{})

	if len(result.Opened) != 0 {
		t.Errorf("expected 0 opened ports, got %d", len(result.Opened))
	}
	if len(result.Closed) != 0 {
		t.Errorf("expected 0 closed ports, got %d", len(result.Closed))
	}
}
