package middleware

import "net/http"

// allowedOrigin is the frontend origin we trust. !TODO: Replace with injecting environment variable
const allowedOrigin = "http://localhost:3000"

// CORSMiddleware serves as middleware to handle CORS in HTTP requests. Returns
// an HTTP handler that adds the Access‑Control‑* headers needed for browsers to
// allow the request.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == allowedOrigin {
			// Only allows requests if the incoming origin is our specified origin.
			w.Header().Set("Access-Control-Allow-Origin", origin)
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
