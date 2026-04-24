package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snapshot.json")
}

func TestStore_SaveAndLoad(t *testing.T) {
	path := tempPath(t)
	store := snapshot.NewStore(path)

	orig := snapshot.Snapshot{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Ports:     []int{22, 80, 443},
	}

	if err := store.Save(orig); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if !got.Timestamp.Equal(orig.Timestamp) {
		t.Errorf("timestamp: got %v, want %v", got.Timestamp, orig.Timestamp)
	}
	if len(got.Ports) != len(orig.Ports) {
		t.Fatalf("ports len: got %d, want %d", len(got.Ports), len(orig.Ports))
	}
	for i, p := range orig.Ports {
		if got.Ports[i] != p {
			t.Errorf("port[%d]: got %d, want %d", i, got.Ports[i], p)
		}
	}
}

func TestStore_LoadMissingFile(t *testing.T) {
	store := snapshot.NewStore("/nonexistent/path/snap.json")
	snap, err := store.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if snap.Ports != nil && len(snap.Ports) != 0 {
		t.Errorf("expected empty ports, got %v", snap.Ports)
	}
}

func TestStore_SaveOverwrites(t *testing.T) {
	path := tempPath(t)
	store := snapshot.NewStore(path)

	_ = store.Save(snapshot.Snapshot{Ports: []int{8080}})
	_ = store.Save(snapshot.Snapshot{Ports: []int{9090}})

	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(got.Ports) != 1 || got.Ports[0] != 9090 {
		t.Errorf("expected [9090], got %v", got.Ports)
	}
}

func TestNew_CopiesPorts(t *testing.T) {
	original := []int{22, 80}
	snap := snapshot.New(original)
	original[0] = 999

	if snap.Ports[0] == 999 {
		t.Error("New should copy the ports slice, not reference it")
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	_ = os.Getenv("") // suppress unused import warning
}
