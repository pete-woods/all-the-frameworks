package api

import (
	"context"
	"net/http"
)

func basicAuth(handler http.HandlerFunc, users map[string]string, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || users[user] != pass {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			_, err := w.Write([]byte("Unauthorised.\n"))
			if err != nil {
				panic(err)
			}
			return
		}

		handler(w, r.WithContext(context.WithValue(r.Context(), "username", user)))
	}
}
