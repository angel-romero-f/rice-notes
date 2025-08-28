package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/angel-romero-f/rice-notes/internal/services"
)

// mockAuthService implements the AuthService interface for testing
type mockAuthService struct {
	authURL             string
	authResult          *services.AuthResult
	authError           error
	validateResult      *services.JWTClaims
	validateError       error
	shouldFailValidation bool
}

func (m *mockAuthService) GetGoogleAuthURL(state string) string {
	return m.authURL
}

func (m *mockAuthService) ExchangeCodeForToken(ctx context.Context, code string) (*services.AuthResult, error) {
	if m.authError != nil {
		return nil, m.authError
	}
	return m.authResult, nil
}

func (m *mockAuthService) ValidateJWT(ctx context.Context, tokenString string) (*services.JWTClaims, error) {
	if m.shouldFailValidation || m.validateError != nil {
		return nil, m.validateError
	}
	return m.validateResult, nil
}

func TestAuthHandler_GoogleLogin(t *testing.T) {
	tests := []struct {
		name           string
		state          string
		expectedURL    string
		expectedStatus int
		expectLocation bool
	}{
		{
			name:           "successful redirect with custom state",
			state:          "custom-state-123",
			expectedURL:    "https://accounts.google.com/oauth/authorize?client_id=test&state=custom-state-123",
			expectedStatus: http.StatusTemporaryRedirect,
			expectLocation: true,
		},
		{
			name:           "successful redirect with auto-generated state",
			state:          "",
			expectedURL:    "https://accounts.google.com/oauth/authorize?client_id=test&state=auto-generated",
			expectedStatus: http.StatusTemporaryRedirect,
			expectLocation: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := &mockAuthService{
				authURL: tt.expectedURL,
			}

			handler := NewAuthHandler(mockService)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/api/auth/google", nil)
			if tt.state != "" {
				q := req.URL.Query()
				q.Add("state", tt.state)
				req.URL.RawQuery = q.Encode()
			}

			rr := httptest.NewRecorder()

			// Execute
			handler.GoogleLogin(rr, req)

			// Assert status
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("GoogleLogin() status = %v, want %v", status, tt.expectedStatus)
			}

			// Assert Location header is set
			if tt.expectLocation {
				location := rr.Header().Get("Location")
				if location == "" {
					t.Error("Expected Location header to be set")
				}
				if !strings.Contains(location, "accounts.google.com") {
					t.Errorf("Expected Location to contain Google OAuth URL, got %v", location)
				}
			}
		})
	}
}

func TestAuthHandler_GoogleCallback(t *testing.T) {
	tests := []struct {
		name               string
		code               string
		state              string
		authResult         *services.AuthResult
		authError          error
		expectedStatus     int
		expectCookie       bool
		expectedRedirectTo string
		expectError        bool
	}{
		{
			name:  "successful callback with valid code",
			code:  "valid-auth-code",
			state: "valid-state",
			authResult: &services.AuthResult{
				Email:   "test@rice.edu",
				Name:    "Test User",
				Picture: "https://example.com/pic.jpg",
				JWT:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			},
			expectedStatus:     http.StatusTemporaryRedirect,
			expectCookie:       true,
			expectedRedirectTo: "http://localhost:3000/dashboard",
		},
		{
			name:           "missing authorization code",
			code:           "",
			state:          "valid-state",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "missing state parameter",
			code:           "valid-code",
			state:          "",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "invalid authorization code",
			code:           "invalid-code",
			state:          "valid-state",
			authError:      errors.New("invalid authorization code"),
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:           "non-rice email",
			code:           "valid-code",
			state:          "valid-state",
			authError:      errors.New("only Rice University emails are allowed"),
			expectedStatus: http.StatusForbidden,
			expectError:    true,
		},
		{
			name:           "unverified email",
			code:           "valid-code",
			state:          "valid-state",
			authError:      errors.New("email not verified"),
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:           "service error",
			code:           "valid-code",
			state:          "valid-state",
			authError:      errors.New("internal service error"),
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := &mockAuthService{
				authResult: tt.authResult,
				authError:  tt.authError,
			}

			handler := NewAuthHandler(mockService)

			// Create request with query parameters
			reqURL := "/api/auth/google/callback"
			if tt.code != "" || tt.state != "" {
				params := url.Values{}
				if tt.code != "" {
					params.Add("code", tt.code)
				}
				if tt.state != "" {
					params.Add("state", tt.state)
				}
				reqURL += "?" + params.Encode()
			}

			req := httptest.NewRequest(http.MethodGet, reqURL, nil)
			rr := httptest.NewRecorder()

			// Execute
			handler.GoogleCallback(rr, req)

			// Assert status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("GoogleCallback() status = %v, want %v", status, tt.expectedStatus)
			}

			// Assert cookie is set for successful auth
			if tt.expectCookie {
				cookies := rr.Result().Cookies()
				found := false
				for _, cookie := range cookies {
					if cookie.Name == "jwt" {
						found = true
						// Verify cookie properties
						if !cookie.HttpOnly {
							t.Error("Expected JWT cookie to be HttpOnly")
						}
						if cookie.Secure != true {
							t.Error("Expected JWT cookie to be Secure")
						}
						if cookie.SameSite != http.SameSiteStrictMode {
							t.Error("Expected JWT cookie to have SameSite=Strict")
						}
						if cookie.Path != "/" {
							t.Error("Expected JWT cookie path to be '/'")
						}
						break
					}
				}
				if !found {
					t.Error("Expected JWT cookie to be set")
				}
			}

			// Assert redirect for successful auth
			if tt.expectedRedirectTo != "" {
				location := rr.Header().Get("Location")
				if location != tt.expectedRedirectTo {
					t.Errorf("Expected redirect to %v, got %v", tt.expectedRedirectTo, location)
				}
			}

			// Assert error response format for failures
			if tt.expectError {
				contentType := rr.Header().Get("Content-Type")
				if !strings.Contains(contentType, "application/json") {
					t.Error("Expected JSON error response")
				}
			}
		})
	}
}

func TestAuthHandler_GoogleCallback_StateValidation(t *testing.T) {
	// Test state validation logic
	t.Run("validates state parameter", func(t *testing.T) {
		mockService := &mockAuthService{}
		handler := NewAuthHandler(mockService)

		// Test with empty state - should fail
		req := httptest.NewRequest(http.MethodGet, "/api/auth/google/callback?code=test", nil)
		rr := httptest.NewRecorder()

		handler.GoogleCallback(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected 400 for missing state, got %d", rr.Code)
		}
	})
}

func TestAuthHandler_CookieSettings(t *testing.T) {
	t.Run("sets secure cookie properties", func(t *testing.T) {
		mockService := &mockAuthService{
			authResult: &services.AuthResult{
				Email:   "test@rice.edu",
				Name:    "Test User",
				Picture: "https://example.com/pic.jpg",
				JWT:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			},
		}

		handler := NewAuthHandler(mockService)
		req := httptest.NewRequest(http.MethodGet, "/api/auth/google/callback?code=test&state=test", nil)
		rr := httptest.NewRecorder()

		handler.GoogleCallback(rr, req)

		if rr.Code != http.StatusTemporaryRedirect {
			t.Fatalf("Expected successful redirect, got %d", rr.Code)
		}

		cookies := rr.Result().Cookies()
		var jwtCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "jwt" {
				jwtCookie = cookie
				break
			}
		}

		if jwtCookie == nil {
			t.Fatal("JWT cookie not found")
		}

		// Verify all security properties
		tests := []struct {
			name     string
			got      interface{}
			expected interface{}
		}{
			{"HttpOnly", jwtCookie.HttpOnly, true},
			{"Secure", jwtCookie.Secure, true},
			{"SameSite", jwtCookie.SameSite, http.SameSiteStrictMode},
			{"Path", jwtCookie.Path, "/"},
			{"Value", jwtCookie.Value, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.got != tt.expected {
					t.Errorf("Cookie %s = %v, want %v", tt.name, tt.got, tt.expected)
				}
			})
		}
	})
}

func TestNewAuthHandler(t *testing.T) {
	mockService := &mockAuthService{}
	handler := NewAuthHandler(mockService)

	if handler == nil {
		t.Error("NewAuthHandler() returned nil")
	}

	if handler.authService != mockService {
		t.Error("NewAuthHandler() did not set authService correctly")
	}
}