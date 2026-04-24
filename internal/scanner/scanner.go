package scanner

import (
	"fmt"
	"net"
	"time"
)

// Port represents an open port detected on the host.
type Port struct {
	Protocol string
	Number   int
	Address  string
}

// String returns a human-readable representation of the port.
func (p Port) String() string {
	return fmt.Sprintf("%s:%d (%s)", p.Address, p.Number, p.Protocol)
}

// Scanner defines the interface for port scanning.
type Scanner interface {
	Scan() ([]Port, error)
}

// TCPScanner scans a range of TCP ports on localhost.
type TCPScanner struct {
	StartPort int
	EndPort   int
	Timeout   time.Duration
}

// NewTCPScanner creates a TCPScanner with sensible defaults.
func NewTCPScanner(start, end int) *TCPScanner {
	return &TCPScanner{
		StartPort: start,
		EndPort:   end,
		Timeout:   500 * time.Millisecond,
	}
}

// Scan probes each port in the configured range and returns open ones.
func (s *TCPScanner) Scan() ([]Port, error) {
	var open []Port

	for port := s.StartPort; port <= s.EndPort; port++ {
		address := fmt.Sprintf("127.0.0.1:%d", port)
		conn, err := net.DialTimeout("tcp", address, s.Timeout)
		if err != nil {
			continue
		}
		conn.Close()
		open = append(open, Port{
			Protocol: "tcp",
			Number:   port,
			Address:  "127.0.0.1",
		})
	}

	return open, nil
}
