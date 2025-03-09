package middleware

import (
	"net/http"
)

// CORS middleware to handle Cross-Origin Resource Sharing
func CORSHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins for simplicity (you should restrict this in production)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		// If it's a pre-flight OPTIONS request, just respond with 200 OK
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue to the next handler (like the actual endpoint handler)
		next.ServeHTTP(w, r)
	})
}
