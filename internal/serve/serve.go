// Package serve previews a generated docs site over HTTP.
package serve

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Handler serves the files under root. A request for a directory is
// resolved to that directory's index.md when one exists; directories
// without an index file return 404 instead of a listing, so the
// preview never exposes more than the generated site itself.
func Handler(root string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rel := path.Clean("/" + strings.TrimPrefix(r.URL.Path, "/"))
		rel = strings.TrimPrefix(rel, "/")
		full := filepath.Join(root, filepath.FromSlash(rel))
		info, err := os.Stat(full)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if info.IsDir() {
			index := filepath.Join(full, "index.md")
			if _, err := os.Stat(index); err != nil {
				http.NotFound(w, r)
				return
			}
			full = index
		}
		http.ServeFile(w, r, full)
	})
}
