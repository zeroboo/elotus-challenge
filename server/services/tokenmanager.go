package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// TokenManager implements ITokenManager using JWT tokens
type TokenManager struct {
	secret                 []byte
	tokenExpirationSeconds int64
}

func NewTokenManager(secret string, tokenExpirationSeconds int64) ITokenManager {
	manager := &TokenManager{
		secret:                 []byte(secret),
		tokenExpirationSeconds: tokenExpirationSeconds,
	}
	return manager
}

// GenerateToken creates a new JWT token for the user using HS256
func (s *TokenManager) GenerateToken(userID int, username string) (string, error) {
	// Create header
	header := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}

	// Create payload with expiration (24 hours)
	now := time.Now().Unix()
	payload := Claims{
		UserID:    userID,
		Username:  username,
		IssuedAt:  now,
		ExpiresAt: now + s.tokenExpirationSeconds,
	}

	// Encode header
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerBytes)

	// Encode payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)

	// Create signature
	message := headerEncoded + "." + payloadEncoded
	signature := s.createSignature(message)

	// Return complete JWT
	return message + "." + signature, nil
}

// ValidateToken validates and parses a JWT token
func (s *TokenManager) ValidateToken(tokenString string) (*Claims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	headerEncoded, payloadEncoded, signatureProvided := parts[0], parts[1], parts[2]

	// Verify signature
	message := headerEncoded + "." + payloadEncoded
	expectedSignature := s.createSignature(message)
	if signatureProvided != expectedSignature {
		return nil, fmt.Errorf("invalid signature")
	}

	// Decode payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return nil, fmt.Errorf("invalid payload encoding")
	}

	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, fmt.Errorf("invalid payload format")
	}

	// Check expiration
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, fmt.Errorf("token has expired")
	}

	return &claims, nil
}

// ExtractTokenFromHeader extracts JWT token from "Bearer <token>" format
// Returns the token, or empty string if extraction fails
func (s *TokenManager) ExtractTokenFromHeader(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

// HasValidBearerFormat checks if the header has valid Bearer format (not necessarily valid token)
func (s *TokenManager) HasValidBearerFormat(authHeader string) bool {
	return len(authHeader) >= 7 && authHeader[:7] == "Bearer "
}

// createSignature creates HMAC-SHA256 signature
func (s *TokenManager) createSignature(message string) string {
	h := hmac.New(sha256.New, s.secret)
	h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
