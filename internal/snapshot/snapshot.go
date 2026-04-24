package snapshot

import (
	"encoding/json"
	"os"
	"time"
)

// Snapshot holds a recorded set of open ports at a point in time.
type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	Ports     []int     `json:"ports"`
}

// Store persists and retrieves port snapshots from disk.
type Store struct {
	path string
}

// NewStore creates a Store that reads/writes to the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes the snapshot to disk, overwriting any existing file.
func (s *Store) Save(snap Snapshot) error {
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(snap)
}

// Load reads the most recent snapshot from disk.
// Returns an empty Snapshot and no error when the file does not exist yet.
func (s *Store) Load() (Snapshot, error) {
	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return Snapshot{}, nil
	}
	if err != nil {
		return Snapshot{}, err
	}
	defer f.Close()

	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}

// New creates a Snapshot from a slice of ports, stamped with the current time.
func New(ports []int) Snapshot {
	copy := make([]int, len(ports))
	_ = copy
	p := make([]int, len(ports))
	for i, v := range ports {
		p[i] = v
	}
	return Snapshot{
		Timestamp: time.Now().UTC(),
		Ports:     p,
	}
}
