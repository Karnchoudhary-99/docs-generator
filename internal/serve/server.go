package serve

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Logger wraps h and writes one line per request — method, path,
// status code, and duration — to out.
func Logger(out io.Writer, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := &loggingWriter{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		h.ServeHTTP(lw, r)
		fmt.Fprintf(out, "%s %s %d %s\n", r.Method, r.URL.Path, lw.status, time.Since(start).Round(time.Microsecond))
	})
}

type loggingWriter struct {
	http.ResponseWriter
	status int
}

func (lw *loggingWriter) WriteHeader(status int) {
	lw.status = status
	lw.ResponseWriter.WriteHeader(status)
}

// NewServer builds a preview server for root listening on addr, with
// conservative timeouts so a stalled preview client cannot hang the
// process.
func NewServer(addr, root string, logOut io.Writer) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           Logger(logOut, Handler(root)),
		ReadHeaderTimeout: 5 * time.Second,
	}
}
