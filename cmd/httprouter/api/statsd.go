package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/felixge/httpsnoop"
)

func statsdMiddleWare(stats *statsd.Client, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrappedWriter := newResponseWriter(w)

		defer func() {
			if wrappedWriter.Status == 0 {
				wrappedWriter.Status = 200
			}

			path := r.URL.Path
			// httprouter doesn't tell you which path matched, or name the routes, so this removes the high card IDs
			if r.Method == "PUT" || r.Method == "DELETE" {
				i := strings.LastIndex(path, "/")
				if i != -1 {
					path = path[0:i]
				}
			}

			tags := []string{
				fmt.Sprintf("request.method:%s", r.Method),
				fmt.Sprintf("request.route:%s", path),
				fmt.Sprintf("response.status_code:%d", wrappedWriter.Status),
			}
			_ = stats.Timing("handler", time.Since(start), tags, 1)
			_ = stats.Incr("handler_count", tags, 1)
		}()
		handler.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	// Wrapped is not embedded to prevent ResponseWriter from directly
	// fulfilling the http.ResponseWriter interface. Wrapping in this
	// way would obscure optional http.ResponseWriter interfaces.
	Wrapped http.ResponseWriter
	Status  int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	var rw responseWriter

	rw.Wrapped = httpsnoop.Wrap(w, httpsnoop.Hooks{
		WriteHeader: func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return func(code int) {
				// The first call to WriteHeader sends the response header.
				// Any subsequent calls are invalid. Only record the first
				// code written.
				if rw.Status == 0 {
					rw.Status = code
				}
				next(code)
			}
		},
	})

	return &rw
}
