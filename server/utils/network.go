package utils

import (
	"net/http"
	"strings"
)

// GetClientIP extracts the client IP address from the request
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (in case of proxy)
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		// Take the first IP in case of multiple IPs
		if idx := strings.Index(forwarded, ","); idx != -1 {
			return strings.TrimSpace(forwarded[:idx])
		}
		return strings.TrimSpace(forwarded)
	}

	// Check X-Real-IP header
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return strings.TrimSpace(realIP)
	}

	// Fall back to RemoteAddr
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx != -1 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}
