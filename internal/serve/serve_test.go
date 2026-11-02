package serve

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "index.md"), []byte("# Home\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "guide"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "guide", "index.md"), []byte("# Guide\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "guide", "deep.md"), []byte("# Deep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "empty"), 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

func get(t *testing.T, h http.Handler, target string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestHandlerServesFiles(t *testing.T) {
	rec := get(t, Handler(testRoot(t)), "/guide/deep.md")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "# Deep") {
		t.Fatalf("body = %q, want guide content", rec.Body.String())
	}
}

func TestHandlerDirectoryResolvesToIndex(t *testing.T) {
	rec := get(t, Handler(testRoot(t)), "/guide/")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "# Guide") {
		t.Fatalf("body = %q, want guide index", rec.Body.String())
	}
}

func TestHandlerDirectoryWithoutIndexIs404NotListing(t *testing.T) {
	rec := get(t, Handler(testRoot(t)), "/empty/")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "deep.md") {
		t.Fatalf("body lists directory contents: %q", rec.Body.String())
	}
}

func TestHandlerMissingFileIs404(t *testing.T) {
	rec := get(t, Handler(testRoot(t)), "/nope.md")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestLoggerRecordsMethodPathAndStatus(t *testing.T) {
	var buf bytes.Buffer
	h := Logger(&buf, Handler(testRoot(t)))
	get(t, h, "/guide/deep.md")
	get(t, h, "/nope.md")
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("logged %d lines, want 2: %q", len(lines), buf.String())
	}
	if !strings.Contains(lines[0], "GET /guide/deep.md 200") {
		t.Fatalf("line 0 = %q, want GET /guide/deep.md 200", lines[0])
	}
	if !strings.Contains(lines[1], "GET /nope.md 404") {
		t.Fatalf("line 1 = %q, want GET /nope.md 404", lines[1])
	}
}

func TestNewServerWiresHandlerAndAddr(t *testing.T) {
	srv := NewServer("127.0.0.1:0", testRoot(t), new(bytes.Buffer))
	if srv.Addr != "127.0.0.1:0" {
		t.Fatalf("Addr = %q", srv.Addr)
	}
	if srv.Handler == nil {
		t.Fatal("Handler is nil")
	}
}
