package test

import (
	"testing"
	"time"

	"elotuschallenge/test/share"
)

func TestGenerateJWT_ValidInput_Success(t *testing.T) {
	jwtService := share.CreateTestJWTService()
	userID := 123
	username := "testuser"

	token, err := jwtService.GenerateToken(userID, username)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token")
	}

	// Token should have 3 parts separated by dots
	parts := len(token)
	if parts < 10 { // Rough check for JWT format
		t.Error("Token appears to be too short for a valid JWT")
	}
}

func TestValidateJWT_ValidToken_Success(t *testing.T) {
	jwtService := share.CreateTestJWTService()
	userID := 456
	username := "validuser"

	// Generate token
	token, err := jwtService.GenerateToken(userID, username)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	// Validate token
	claims, err := jwtService.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	// Check claims
	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}

	if claims.Username != username {
		t.Errorf("Expected Username %s, got %s", username, claims.Username)
	}

	// Check expiration is in the future
	if claims.ExpiresAt <= time.Now().Unix() {
		t.Error("Token should not be expired")
	}

	// Check issued time is recent
	if claims.IssuedAt > time.Now().Unix() {
		t.Error("Issued time should not be in the future")
	}
}

func TestValidateJWT_InvalidToken_Error(t *testing.T) {
	jwtService := share.CreateTestJWTService()
	testCases := []struct {
		name  string
		token string
	}{
		{
			name:  "Empty token",
			token: "",
		},
		{
			name:  "Invalid format",
			token: "invalid.token",
		},
		{
			name:  "Random string",
			token: "this.is.not.a.jwt",
		},
		{
			name:  "Too few parts",
			token: "header.payload",
		},
		{
			name:  "Too many parts",
			token: "header.payload.signature.extra",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := jwtService.ValidateToken(tc.token)
			if err == nil {
				t.Error("Expected validation to fail for invalid token")
			}
		})
	}
}

func TestValidateJWT_TamperedToken_Error(t *testing.T) {
	jwtService := share.CreateTestJWTService()
	userID := 789
	username := "tampereduser"

	// Generate valid token
	token, err := jwtService.GenerateToken(userID, username)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	// Tamper with the token by changing one character
	tamperedToken := token[:len(token)-1] + "X"

	// Try to validate tampered token
	_, err = jwtService.ValidateToken(tamperedToken)
	if err == nil {
		t.Error("Expected validation to fail for tampered token")
	}
}

func TestExtractTokenFromBearer_ValidFormat_Success(t *testing.T) {
	jwtService := share.CreateTestJWTService()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.example.token"
	authHeader := "Bearer " + token

	extracted := jwtService.ExtractTokenFromHeader(authHeader)
	if extracted != token {
		t.Errorf("Expected token '%s', got '%s'", token, extracted)
	}
}

func TestExtractTokenFromBearer_InvalidFormat_Empty(t *testing.T) {
	jwtService := share.CreateTestJWTService()
	testCases := []struct {
		name   string
		header string
	}{
		{
			name:   "Empty header",
			header: "",
		},
		{
			name:   "No Bearer prefix",
			header: "token",
		},
		{
			name:   "Wrong prefix",
			header: "Basic token",
		},
		{
			name:   "Bearer without space",
			header: "Bearertoken",
		},
		{
			name:   "Bearer with no token",
			header: "Bearer ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			extracted := jwtService.ExtractTokenFromHeader(tc.header)
			if extracted != "" {
				t.Errorf("Expected empty token for invalid format, got '%s'", extracted)
			}
		})
	}
}

func TestJWT_GenerateAndValidate_RoundTrip(t *testing.T) {
	jwtService := share.CreateTestJWTService()
	testCases := []struct {
		userID   int
		username string
	}{
		{1, "user1"},
		{999, "longusernametest"},
		{0, "zerouser"},
		{-1, "negativeuser"}, // Edge case
	}

	for _, tc := range testCases {
		t.Run(tc.username, func(t *testing.T) {
			// Generate token
			token, err := jwtService.GenerateToken(tc.userID, tc.username)
			if err != nil {
				t.Fatalf("Failed to generate JWT: %v", err)
			}

			// Validate token
			claims, err := jwtService.ValidateToken(token)
			if err != nil {
				t.Fatalf("Failed to validate JWT: %v", err)
			}

			// Verify claims match
			if claims.UserID != tc.userID {
				t.Errorf("UserID mismatch: expected %d, got %d", tc.userID, claims.UserID)
			}

			if claims.Username != tc.username {
				t.Errorf("Username mismatch: expected %s, got %s", tc.username, claims.Username)
			}
		})
	}
}
