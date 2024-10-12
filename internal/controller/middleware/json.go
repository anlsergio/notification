package middleware

import "net/http"

// SetJSONContent adds the necessary headers to serve JSON content.
func SetJSONContent(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next(w, r)
	}
}
