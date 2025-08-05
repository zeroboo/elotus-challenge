package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"elotuschallenge/internal"
	"elotuschallenge/transfer"

	"github.com/rs/zerolog/log"
)

const ErrorUnauthorized = "Unauthorized access"

// AuthUser validates JWT tokens for protected routes
func AuthUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			sendUnauthorized(w, "Authorization header required")
			return
		}

		// Check if header has valid Bearer format
		if !internal.TokenManager.HasValidBearerFormat(authHeader) {
			sendUnauthorized(w, "Invalid authorization format")
			return
		}

		// Extract token from Bearer format
		token := internal.TokenManager.ExtractTokenFromHeader(authHeader)
		if token == "" {
			sendUnauthorized(w, "Invalid or expired token")
			return
		}

		// Validate JWT token
		claims, err := internal.TokenManager.ValidateToken(token)
		if err != nil {
			log.Error().Err(err).Msg("JWT validation failed")
			sendUnauthorized(w, "Invalid or expired token")
			return
		}

		// Add user info to request context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		r = r.WithContext(ctx)

		// Call next handler
		next(w, r)
	}
}

// sendUnauthorized sends a uniform unauthorized response
func sendUnauthorized(w http.ResponseWriter, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	response := transfer.NewErrorResponse(ErrorUnauthorized)
	json.NewEncoder(w).Encode(response)
}
