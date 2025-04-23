package middleware

import "net/http"

// CorsMiddleware sets response headers for all requests
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// If preflight request, respond with headers and 200
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}
