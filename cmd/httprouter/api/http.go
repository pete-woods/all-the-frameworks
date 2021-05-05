package api

import (
	"encoding/json"
	"net/http"
)

func writeError(w http.ResponseWriter, code int, err interface{}) {
	if err == nil {
		return
	}

	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return
	}

	msg := ""
	if e, ok := err.(error); ok {
		msg = e.Error()
	} else if s, ok := err.(string); ok {
		msg = s
	}

	type httpError struct {
		Error string `json:"error"`
	}
	body := httpError{
		Error: msg,
	}

	writeJSON(w, code, body)
}

func writeJSON(w http.ResponseWriter, code int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if resp != nil {
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			panic(err)
		}
	}
}

func panicHandler(w http.ResponseWriter, _ *http.Request, i interface{}) {
	writeError(w, http.StatusInternalServerError, i)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	writeError(w, http.StatusNotFound, "routing error")
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
}
