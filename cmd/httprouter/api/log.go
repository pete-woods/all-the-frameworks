package api

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func logMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrappedWriter := newResponseWriter(w)

		defer func() {
			if wrappedWriter.Status == 0 {
				wrappedWriter.Status = 200
			}

			log.WithFields(log.Fields{
				"method":      r.Method,
				"url":         r.URL.Path,
				"status_code": wrappedWriter.Status,
				"duration":    time.Since(start),
			}).Info()
		}()
		handler.ServeHTTP(w, r)
	})
}
