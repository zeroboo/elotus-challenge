package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"elotuschallenge/internal"
	"elotuschallenge/share"
	"elotuschallenge/transfer"
	"elotuschallenge/utils"

	"github.com/rs/zerolog/log"
)

const ErrMsgUnauthorized = "Unauthorized access"

var ErrNoAuthorizationHeader = fmt.Errorf("authorization header required")
var ErrInvalidAuthorizationFormat = fmt.Errorf("invalid authorization format")
var ErrMalformedToken = fmt.Errorf("malformed token")
var ErrInvalidToken = fmt.Errorf("invalid token")

// AuthUser validates JWT tokens for protected routes
func AuthUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get(share.HeaderAuthorization)
		if authHeader == "" {
			ResponseUnauthorized(w, r, ErrNoAuthorizationHeader)
			return
		}

		// Check if header has valid Bearer format
		if !internal.TokenManager.HasValidBearerFormat(authHeader) {
			ResponseUnauthorized(w, r, ErrInvalidAuthorizationFormat)
			return
		}

		// Extract token from Bearer format
		token := internal.TokenManager.ExtractTokenFromHeader(authHeader)
		if token == "" {
			ResponseUnauthorized(w, r, ErrMalformedToken)
			return
		}

		// Validate JWT token
		claims, errValidate := internal.TokenManager.ValidateToken(token)
		if errValidate != nil {
			ResponseUnauthorized(w, r, fmt.Errorf("%w: %w", ErrInvalidToken, errValidate))
			return
		}

		// Add user info to request context
		ctx := context.WithValue(r.Context(), share.ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, share.ContextKeyUsername, claims.Username)
		r = r.WithContext(ctx)

		// Call next handler
		next(w, r)
	}
}

// ResponseUnauthorized sends a uniform unauthorized response
func ResponseUnauthorized(resp http.ResponseWriter, req *http.Request, authorizeError error) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusUnauthorized)

	response := transfer.NewErrorResponse(ErrMsgUnauthorized)
	json.NewEncoder(resp).Encode(response)
	var logMsg string
	if authorizeError != nil {
		logMsg = authorizeError.Error()
	}

	log.Error().
		Str("client_ip", utils.GetClientIP(req)).
		Str("error", logMsg).Msg("Unauthorized request")
}
