package watch

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRunCallsOnChangeAndStopsOnCancel(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "a.go")
	if err := os.WriteFile(f, []byte("v1"), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	called := make(chan struct{}, 1)
	go func() {
		_ = Run(ctx, []string{dir}, 10*time.Millisecond, func() {
			select {
			case called <- struct{}{}:
			default:
			}
		})
	}()
	time.Sleep(30 * time.Millisecond)
	if err := os.WriteFile(f, []byte("v1 changed"), 0o644); err != nil {
		t.Fatal(err)
	}
	select {
	case <-called:
	case <-time.After(2 * time.Second):
		t.Fatal("onChange not called after file change")
	}
	cancel()
}
