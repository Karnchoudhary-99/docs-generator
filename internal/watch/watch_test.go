package watch

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSnapshotOnlyTracksGoAndMarkdown(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.go", "b.md", "c.txt", "d.json"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	snap := Snapshot(dir)
	if len(snap) != 2 {
		t.Fatalf("snapshot has %d files, want 2 (.go + .md)", len(snap))
	}
}

func TestChangedDetectsModifyAddRemove(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "a.go")
	if err := os.WriteFile(f, []byte("one"), 0o644); err != nil {
		t.Fatal(err)
	}
	before := Snapshot(dir)
	if Changed(before, before) {
		t.Fatal("identical snapshots reported as changed")
	}
	// Modify: size change.
	if err := os.WriteFile(f, []byte("one two"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !Changed(before, Snapshot(dir)) {
		t.Fatal("modified file not detected")
	}
	// Add.
	mid := Snapshot(dir)
	if err := os.WriteFile(filepath.Join(dir, "b.md"), []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !Changed(mid, Snapshot(dir)) {
		t.Fatal("added file not detected")
	}
	// Remove.
	afterAdd := Snapshot(dir)
	if err := os.Remove(f); err != nil {
		t.Fatal(err)
	}
	if !Changed(afterAdd, Snapshot(dir)) {
		t.Fatal("removed file not detected")
	}
}
