package middleware

import (
	"net/http"
	"os"
	"strings"
)

// getAllowedOrigins returns the list of allowed origins for CORS
func getAllowedOrigins() []string {
	origins := []string{"http://localhost:3000"} // Always allow local development
	
	// Add production origins from environment variable
	if prodOrigins := os.Getenv("ALLOWED_ORIGINS"); prodOrigins != "" {
		origins = append(origins, strings.Split(prodOrigins, ",")...)
	}
	
	return origins
}

// CORSMiddleware serves as middleware to handle CORS in HTTP requests. Returns
// an HTTP handler that adds the Access‑Control‑* headers needed for browsers to
// allow the request.
func CORSMiddleware(next http.Handler) http.Handler {
	allowedOrigins := getAllowedOrigins()
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		
		// Check if origin is in allowed list
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		// Which HTTP methods are permitted in cross‑origin requests.
		w.Header().Set("Access-Control-Allow-Methods",
			"GET, POST, PUT, PATCH, DELETE, OPTIONS")

		// Which request headers the browser may send.
		w.Header().Set("Access-Control-Allow-Headers",
			"Authorization, Content-Type")

		// Tell the browser to include cookies / authorization headers.
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// Passes the request to the next handler to be ran in the middleware chain.
		next.ServeHTTP(w, r)
	})
}
