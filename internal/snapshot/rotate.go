package snapshot

import (
	"fmt"
	"os"
	"path/filepath"	
	"time"
)

// RotatingStore wraps Store and keeps up to MaxBackups previous snapshots
// alongside the current one, named with a timestamp suffix.
type RotatingStore struct {
	base       string
	maxBackups int
}

// NewRotatingStore creates a RotatingStore. base is the primary snapshot file
// path; backups are written next to it with a timestamp suffix.
func NewRotatingStore(base string, maxBackups int) *RotatingStore {
	if maxBackups < 1 {
		maxBackups = 3
	}
	return &RotatingStore{base: base, maxBackups: maxBackups}
}

// Save rotates existing snapshot to a timestamped backup, then writes snap.
func (r *RotatingStore) Save(snap Snapshot) error {
	if _, err := os.Stat(r.base); err == nil {
		if err := r.rotate(); err != nil {
			return fmt.Errorf("rotate: %w", err)
		}
	}
	return NewStore(r.base).Save(snap)
}

// Load returns the most recent snapshot from the primary file.
func (r *RotatingStore) Load() (Snapshot, error) {
	return NewStore(r.base).Load()
}

func (r *RotatingStore) rotate() error {
	dir := filepath.Dir(r.base)
	ext := filepath.Ext(r.base)
	name := filepath.Base(r.base[:len(r.base)-len(ext)])

	ts := time.Now().UTC().Format("20060102T150405")
	dst := filepath.Join(dir, fmt.Sprintf("%s.%s%s", name, ts, ext))

	if err := os.Rename(r.base, dst); err != nil {
		return err
	}

	return r.pruneBackups(dir, name, ext)
}

func (r *RotatingStore) pruneBackups(dir, name, ext string) error {
	pattern := filepath.Join(dir, fmt.Sprintf("%s.*%s", name, ext))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for len(matches) > r.maxBackups {
		if err := os.Remove(matches[0]); err != nil {
			return err
		}
		matches = matches[1:]
	}
	return nil
}
