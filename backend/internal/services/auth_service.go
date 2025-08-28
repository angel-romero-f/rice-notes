package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// OAuth2Provider defines the interface for OAuth2 operations
type OAuth2Provider interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*TokenResult, error)
	GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error)
}

// TokenResult represents the result of token exchange
type TokenResult struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// UserInfo represents Google user information
type UserInfo struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Verified bool   `json:"email_verified"`
}

// AuthResult represents the result of successful authentication
type AuthResult struct {
	Email   string
	Name    string
	Picture string
	JWT     string
}

// JWTClaims represents the claims in our JWT
type JWTClaims struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	jwt.RegisteredClaims
}

// GoogleOAuth2Provider implements OAuth2Provider for Google
type GoogleOAuth2Provider struct {
	config *oauth2.Config
}

// NewGoogleOAuth2Provider creates a new Google OAuth2 provider
func NewGoogleOAuth2Provider(clientID, clientSecret, redirectURL string) *GoogleOAuth2Provider {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}

	return &GoogleOAuth2Provider{config: config}
}

// GetAuthURL generates the Google OAuth2 authorization URL
func (g *GoogleOAuth2Provider) GetAuthURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCode exchanges authorization code for access token
func (g *GoogleOAuth2Provider) ExchangeCode(ctx context.Context, code string) (*TokenResult, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		slog.Error("Failed to exchange code for token", "error", err)
		return nil, fmt.Errorf("code exchange failed: %w", err)
	}

	return &TokenResult{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresIn:   int(time.Until(token.Expiry).Seconds()),
	}, nil
}

// GetUserInfo retrieves user information from Google
func (g *GoogleOAuth2Provider) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	url := "https://www.googleapis.com/oauth2/v2/userinfo"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to get user info", "error", err)
		return nil, fmt.Errorf("userinfo request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Userinfo request returned non-200 status", "status", resp.StatusCode)
		return nil, fmt.Errorf("userinfo request failed with status: %d", resp.StatusCode)
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		slog.Error("Failed to decode user info", "error", err)
		return nil, fmt.Errorf("failed to decode userinfo: %w", err)
	}

	slog.Debug("Retrieved user info", "email", userInfo.Email, "verified", userInfo.Verified)
	return &userInfo, nil
}

// AuthService handles authentication operations
type AuthService struct {
	provider  OAuth2Provider
	jwtSecret []byte
}

// NewAuthService creates a new AuthService instance
func NewAuthService(provider OAuth2Provider, jwtSecret string) *AuthService {
	return &AuthService{
		provider:  provider,
		jwtSecret: []byte(jwtSecret),
	}
}

// GetGoogleAuthURL generates a Google OAuth2 authorization URL with state
func (a *AuthService) GetGoogleAuthURL(state string) string {
	if state == "" {
		// Generate a random state if none provided
		state = a.generateState()
	}

	url := a.provider.GetAuthURL(state)
	slog.Info("Generated Google auth URL", "state", state)
	return url
}

// ExchangeCodeForToken exchanges an authorization code for a JWT token
func (a *AuthService) ExchangeCodeForToken(ctx context.Context, code string) (*AuthResult, error) {
	slog.Info("Starting code exchange", "code_length", len(code))

	// Exchange code for access token
	tokenResult, err := a.provider.ExchangeCode(ctx, code)
	if err != nil {
		slog.Error("Code exchange failed", "error", err)
		return nil, fmt.Errorf("code exchange failed: %w", err)
	}

	// Get user information
	userInfo, err := a.provider.GetUserInfo(ctx, tokenResult.AccessToken)
	if err != nil {
		slog.Error("Failed to get user info", "error", err)
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Validate Rice University email first
	if !a.isRiceEmail(userInfo.Email) {
		slog.Warn("Non-Rice email attempted login", "email", userInfo.Email)
		return nil, errors.New("only Rice University emails are allowed")
	}

	// For Rice emails, we trust Google's domain verification
	// For non-Rice emails (if we ever allow them), require email verification
	if !a.isRiceEmail(userInfo.Email) && !userInfo.Verified {
		slog.Warn("User email not verified", "email", userInfo.Email)
		return nil, errors.New("email not verified")
	}

	// Generate JWT
	jwtToken, err := a.generateJWT(userInfo)
	if err != nil {
		slog.Error("Failed to generate JWT", "error", err)
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	slog.Info("Successful authentication", "email", userInfo.Email)

	return &AuthResult{
		Email:   userInfo.Email,
		Name:    userInfo.Name,
		Picture: userInfo.Picture,
		JWT:     jwtToken,
	}, nil
}

// ValidateJWT validates a JWT token and returns claims
func (a *AuthService) ValidateJWT(ctx context.Context, tokenString string) (*JWTClaims, error) {
	if tokenString == "" {
		return nil, errors.New("empty token")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.jwtSecret, nil
	})

	if err != nil {
		slog.Error("JWT validation failed", "error", err)
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		slog.Error("Invalid JWT claims")
		return nil, errors.New("invalid token claims")
	}

	// Check if token is expired
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		slog.Warn("Expired JWT token", "email", claims.Email)
		return nil, errors.New("token expired")
	}

	slog.Debug("JWT validation successful", "email", claims.Email)
	return claims, nil
}

// isRiceEmail checks if an email belongs to Rice University
func (a *AuthService) isRiceEmail(email string) bool {
	if email == "" {
		return false
	}

	// Convert to lowercase for case-insensitive comparison
	email = strings.ToLower(email)

	// Check for @rice.edu or @subdomain.rice.edu
	return strings.HasSuffix(email, "@rice.edu") || strings.Contains(email, ".rice.edu")
}

// generateJWT creates a JWT token for the authenticated user
func (a *AuthService) generateJWT(userInfo *UserInfo) (string, error) {
	// Token expires in 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &JWTClaims{
		Email:   userInfo.Email,
		Name:    userInfo.Name,
		Picture: userInfo.Picture,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "rice-notes",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	slog.Debug("Generated JWT token", "email", userInfo.Email, "expires", expirationTime)
	return tokenString, nil
}

// generateState generates a random state string for OAuth2 flow
func (a *AuthService) generateState() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based state if random generation fails
		return fmt.Sprintf("state_%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}
