package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/angel-romero-f/rice-notes/internal/services"
)

// UserContextKey is the key used to store user information in request context
type UserContextKey string

const userContextKey UserContextKey = "user"

// JWTMiddleware validates JWT tokens from Authorization header or cookies and adds user context
func JWTMiddleware(authService AuthServiceInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to extract JWT from Authorization header first (for cross-origin)
			var tokenString string
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
				slog.Debug("JWT found in Authorization header", "path", r.URL.Path)
			} else {
				// Fallback to cookie (for same-origin)
				cookie, err := r.Cookie("jwt")
				if err != nil {
					slog.Debug("No JWT found in Authorization header or cookie", "path", r.URL.Path)
					http.Error(w, "Authentication required", http.StatusUnauthorized)
					return
				}
				tokenString = cookie.Value
				slog.Debug("JWT found in cookie", "path", r.URL.Path)
			}

			// Validate JWT
			claims, err := authService.ValidateJWT(r.Context(), tokenString)
			if err != nil {
				slog.Warn("Invalid JWT token", "error", err, "path", r.URL.Path)
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Add user information to request context
			ctx := context.WithValue(r.Context(), userContextKey, claims)
			
			slog.Debug("JWT validated successfully", "email", claims.Email, "path", r.URL.Path)
			
			// Continue to the next handler with the updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext extracts user claims from request context
func GetUserFromContext(ctx context.Context) (*services.JWTClaims, bool) {
	user, ok := ctx.Value(userContextKey).(*services.JWTClaims)
	return user, ok
}

// AuthServiceInterface defines the methods needed by the JWT middleware
type AuthServiceInterface interface {
	ValidateJWT(ctx context.Context, tokenString string) (*services.JWTClaims, error)
}