package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/angel-romero-f/rice-notes/internal/services"
)

// AuthService defines the business logic for authentication operations
type AuthService interface {
	GetGoogleAuthURL(state string) string
	ExchangeCodeForToken(ctx context.Context, code string) (*services.AuthResult, error)
	ValidateJWT(ctx context.Context, tokenString string) (*services.JWTClaims, error)
}

// AuthHandler handles HTTP requests for authentication operations
type AuthHandler struct {
	authService AuthService
}

// NewAuthHandler returns a new AuthHandler instance with the provided AuthService
func NewAuthHandler(s AuthService) *AuthHandler {
	return &AuthHandler{authService: s}
}

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// UserResponse represents a user information response
type UserResponse struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// GoogleLogin initiates the Google OAuth2 flow by redirecting to Google's authorization URL
func (a *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	slog.Info("Google login initiated", "remote_addr", r.RemoteAddr, "user_agent", r.UserAgent())

	// Get state parameter from query string or generate one
	state := r.URL.Query().Get("state")
	if state == "" {
		// Generate a simple state for this session
		state = "auth_" + r.Header.Get("X-Request-ID") // Use request ID if available
		if state == "auth_" {
			// Fallback state generation
			state = "auth_request"
		}
	}

	// Get Google OAuth URL from service
	authURL := a.authService.GetGoogleAuthURL(state)
	
	slog.Info("Redirecting to Google OAuth", "state", state, "url_length", len(authURL))

	// Redirect to Google OAuth URL
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// GoogleCallback handles the OAuth2 callback from Google
func (a *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	slog.Info("Google callback received", "remote_addr", r.RemoteAddr)

	// Extract query parameters
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	errorParam := r.URL.Query().Get("error")

	// Check if user denied access
	if errorParam != "" {
		slog.Warn("User denied OAuth access", "error", errorParam)
		a.sendErrorResponse(w, http.StatusUnauthorized, "access_denied", "User denied access")
		return
	}

	// Validate required parameters
	if code == "" {
		slog.Warn("Missing authorization code in callback")
		a.sendErrorResponse(w, http.StatusBadRequest, "missing_code", "Authorization code is required")
		return
	}

	if state == "" {
		slog.Warn("Missing state parameter in callback")
		a.sendErrorResponse(w, http.StatusBadRequest, "missing_state", "State parameter is required")
		return
	}

	// Exchange code for JWT token
	authResult, err := a.authService.ExchangeCodeForToken(r.Context(), code)
	if err != nil {
		slog.Error("Code exchange failed", "error", err, "code_length", len(code))
		
		// Map service errors to appropriate HTTP status codes
		errMsg := err.Error()
		switch {
		case errMsg == "only Rice University emails are allowed":
			a.sendErrorResponse(w, http.StatusForbidden, "non_rice_email", "Only Rice University email addresses are allowed")
		case errMsg == "email not verified":
			a.sendErrorResponse(w, http.StatusUnauthorized, "unverified_email", "Email address must be verified")
		case errMsg == "invalid authorization code":
			a.sendErrorResponse(w, http.StatusUnauthorized, "invalid_code", "Invalid authorization code")
		default:
			a.sendErrorResponse(w, http.StatusInternalServerError, "auth_error", "Authentication failed")
		}
		return
	}

	// Set JWT in secure HttpOnly cookie
	a.setJWTCookie(w, authResult.JWT)

	slog.Info("Successful authentication", "email", authResult.Email, "name", authResult.Name)

	// Redirect to frontend dashboard
	frontendURL := "http://localhost:3000/dashboard" // TODO: Make configurable via environment
	http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
}

// setJWTCookie sets a secure HttpOnly cookie with the JWT token
func (a *AuthHandler) setJWTCookie(w http.ResponseWriter, jwt string) {
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    jwt,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Requires HTTPS in production
		SameSite: http.SameSiteStrictMode,
		MaxAge:   24 * 60 * 60, // 24 hours (matching JWT expiration)
	}

	http.SetCookie(w, cookie)
	slog.Debug("JWT cookie set", "cookie_name", cookie.Name, "max_age", cookie.MaxAge)
}

// sendErrorResponse sends a JSON error response with the specified status code
func (a *AuthHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, errorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Error:   errorCode,
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to encode error response", "error", err)
		// Fallback to plain text if JSON encoding fails
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	slog.Debug("Error response sent", "status", statusCode, "error_code", errorCode, "message", message)
}

// Me returns the current user's information from the JWT token
func (a *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	slog.Info("User info requested", "remote_addr", r.RemoteAddr)

	// Extract JWT from cookie
	cookie, err := r.Cookie("jwt")
	if err != nil {
		slog.Warn("No JWT cookie found")
		a.sendErrorResponse(w, http.StatusUnauthorized, "no_token", "Authentication required")
		return
	}

	// Validate JWT
	claims, err := a.authService.ValidateJWT(r.Context(), cookie.Value)
	if err != nil {
		slog.Warn("Invalid JWT token", "error", err)
		a.sendErrorResponse(w, http.StatusUnauthorized, "invalid_token", "Invalid or expired token")
		return
	}

	// Return user information
	response := UserResponse{
		Email:   claims.Email,
		Name:    claims.Name,
		Picture: claims.Picture,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to encode user response", "error", err)
		a.sendErrorResponse(w, http.StatusInternalServerError, "encoding_error", "Failed to encode response")
		return
	}

	slog.Info("User info returned successfully", "email", claims.Email)
}