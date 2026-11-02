package watch

import (
	"context"
	"time"
)

// Run polls dirs every interval and calls onChange whenever the
// snapshot changes. It returns when ctx is cancelled; onChange runs
// on the polling goroutine, so a slow rebuild delays the next poll
// rather than stacking up.
func Run(ctx context.Context, dirs []string, interval time.Duration, onChange func()) error {
	prev := Snapshot(dirs...)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			next := Snapshot(dirs...)
			if Changed(prev, next) {
				prev = next
				onChange()
			}
		}
	}
}
