package main

import "net/http"

// Readiness
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)                           // -- It can be also w.WriteHeader(200)
	w.Write([]byte(http.StatusText(http.StatusOK) + "\n")) // -- It can be also w.Write([]byte("OK"))
}
