// Package watch rebuilds a docs site when its sources change. It polls
// directory snapshots instead of using OS file events, so it works the
// same on every platform with no extra dependencies.
package watch

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// fileState records what a rebuild depends on for one source file.
type fileState struct {
	size    int64
	modTime time.Time
}

// Snapshot maps every .go and .md file under dirs (recursively) to
// its state. Missing directories contribute nothing; unreadable files
// are skipped, so a half-written tree never breaks a snapshot.
func Snapshot(dirs ...string) map[string]fileState {
	snap := make(map[string]fileState)
	for _, dir := range dirs {
		_ = filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			if !strings.HasSuffix(p, ".go") && !strings.HasSuffix(p, ".md") {
				return nil
			}
			info, err := d.Info()
			if err != nil {
				return nil
			}
			snap[p] = fileState{size: info.Size(), modTime: info.ModTime()}
			return nil
		})
	}
	return snap
}
