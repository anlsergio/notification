package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger decorates HTTP requests with logs.
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("new HTTP request started")
		start := time.Now()
		defer func() {
			log.Printf("HTTP request ended after %s", time.Since(start))
		}()

		next.ServeHTTP(w, r)
	}
}
