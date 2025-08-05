package middlewares

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// AuthUser validates JWT tokens for protected routes
func AuthUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement JWT token validation
		// For now, just pass through
		log.Info().Msg("Auth middleware - TODO: validate JWT token")
		next(w, r)
	}
}
