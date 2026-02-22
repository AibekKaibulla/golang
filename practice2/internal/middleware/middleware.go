package middleware

import (
	"log"
	"net/http"
	"time"
)

const validAPIKey = "test12131p1p31ploxabcdefgjasfnjafnasjfnfjs"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-KEY")
		if key != validAPIKey {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte(`{"error":"unauthorized: invalid or missing X-API-KEY"}`))
            return
        }
        next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("[%s] %s %s", start.Format(time.RFC3339), r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}