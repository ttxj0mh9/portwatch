package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func TestRotatingStore_BasicSaveLoad(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "snap.json")
	rs := snapshot.NewRotatingStore(base, 3)

	snap := snapshot.New([]int{22, 443})
	if err := rs.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := rs.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(got.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(got.Ports))
	}
}

func TestRotatingStore_CreatesBackup(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "snap.json")
	rs := snapshot.NewRotatingStore(base, 5)

	_ = rs.Save(snapshot.New([]int{80}))
	_ = rs.Save(snapshot.New([]int{443}))

	matches, err := filepath.Glob(filepath.Join(dir, "snap.*.json"))
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	if len(matches) != 1 {
		t.Errorf("expected 1 backup file, got %d", len(matches))
	}
}

func TestRotatingStore_PrunesOldBackups(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "snap.json")
	maxBackups := 2
	rs := snapshot.NewRotatingStore(base, maxBackups)

	// Write more snapshots than maxBackups allows
	for i := 0; i < maxBackups+3; i++ {
		// Small sleep would normally differentiate timestamps; skip in tests.
		if err := rs.Save(snapshot.New([]int{i})); err != nil {
			t.Fatalf("Save iteration %d: %v", i, err)
		}
	}

	matches, err := filepath.Glob(filepath.Join(dir, "snap.*.json"))
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	if len(matches) > maxBackups {
		t.Errorf("expected at most %d backups, got %d", maxBackups, len(matches))
	}
}

func TestRotatingStore_LoadMissing(t *testing.T) {
	rs := snapshot.NewRotatingStore("/no/such/file.json", 3)
	snap, err := rs.Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(snap.Ports) != 0 {
		t.Errorf("expected empty ports, got %v", snap.Ports)
	}
	_ = os.Getenv("") // suppress unused import
}
