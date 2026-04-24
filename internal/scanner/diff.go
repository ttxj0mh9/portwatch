package scanner

import "fmt"

// ChangeKind describes whether a port was opened or closed.
type ChangeKind string

const (
	Opened ChangeKind = "opened"
	Closed ChangeKind = "closed"
)

// Change represents a single detected port state change.
type Change struct {
	Kind ChangeKind
	Port Port
}

// String returns a human-readable description of the change.
func (c Change) String() string {
	return fmt.Sprintf("port %s: %s", c.Kind, c.Port)
}

// Diff compares two snapshots of open ports and returns what changed.
// previous and current are slices returned by successive Scan calls.
func Diff(previous, current []Port) []Change {
	prevSet := toSet(previous)
	currSet := toSet(current)

	var changes []Change

	// Ports in current but not in previous → newly opened.
	for key, p := range currSet {
		if _, exists := prevSet[key]; !exists {
			changes = append(changes, Change{Kind: Opened, Port: p})
		}
	}

	// Ports in previous but not in current → closed.
	for key, p := range prevSet {
		if _, exists := currSet[key]; !exists {
			changes = append(changes, Change{Kind: Closed, Port: p})
		}
	}

	return changes
}

// toSet converts a slice of Ports into a map keyed by "protocol:address:number".
func toSet(ports []Port) map[string]Port {
	m := make(map[string]Port, len(ports))
	for _, p := range ports {
		key := fmt.Sprintf("%s:%s:%d", p.Protocol, p.Address, p.Number)
		m[key] = p
	}
	return m
}
