package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"elotuschallenge/common"
	"elotuschallenge/internal"
	"elotuschallenge/middleware"
	"elotuschallenge/transfer"
)

func init() {
	// Initialize services for testing
	internal.InitServices()
}

func TestAuthMiddleware_ValidToken_Success(t *testing.T) {
	// Use the global JWT service that matches the middleware
	userID := 123
	username := "authuser"

	// Generate valid token using the standalone JWT service
	token, err := internal.TokenManager.GenerateToken(userID, username)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	// Create a test handler that checks context
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user info is in context
		contextUserID, ok := r.Context().Value("user_id").(int)
		if !ok {
			t.Error("Expected user_id in context")
			return
		}
		if contextUserID != userID {
			t.Errorf("Expected user_id %d, got %d", userID, contextUserID)
		}

		contextUsername, ok := r.Context().Value("username").(string)
		if !ok {
			t.Error("Expected username in context")
			return
		}
		if contextUsername != username {
			t.Errorf("Expected username %s, got %s", username, contextUsername)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Wrap with auth middleware
	authHandler := middleware.AuthUser(testHandler)

	// Create request with valid token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set(common.HeaderAuthorization, "Bearer "+token)
	w := httptest.NewRecorder()

	// Call middleware
	authHandler(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != "success" {
		t.Errorf("Expected 'success', got '%s'", w.Body.String())
	}
}

func TestAuthMiddleware_NoAuthHeader_Error(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called without auth")
	})

	authHandler := middleware.AuthUser(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	authHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var response transfer.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Message != middleware.ErrMsgUnauthorized {
		t.Errorf("Expected message '%v', got '%s'", middleware.ErrMsgUnauthorized, response.Message)
	}
}

func TestAuthMiddleware_InvalidAuthFormat_Error(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with invalid auth format")
	})

	authHandler := middleware.AuthUser(testHandler)

	testCases := []struct {
		name   string
		header string
	}{
		{"No Bearer prefix", "token123"},
		{"Wrong prefix", "Basic token123"},
		{"Bearer without space", "Bearertoken123"},
		{"Empty Bearer", "Bearer "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set(common.HeaderAuthorization, tc.header)
			w := httptest.NewRecorder()

			authHandler(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
			}

			var response transfer.APIResponse
			json.NewDecoder(w.Body).Decode(&response)

			if response.Success {
				t.Error("Expected success to be false")
			}

			if response.Message != middleware.ErrMsgUnauthorized {
				t.Errorf("Expected message '%v', got '%s'", middleware.ErrMsgUnauthorized, response.Message)
			}
		})
	}
}

func TestAuthMiddleware_InvalidToken_Error(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with invalid token")
	})

	authHandler := middleware.AuthUser(testHandler)

	testCases := []struct {
		name  string
		token string
	}{
		{"Invalid format", "invalid.token"},
		{"Tampered token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.tampered.signature"},
		{"Empty token", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set(common.HeaderAuthorization, "Bearer "+tc.token)
			w := httptest.NewRecorder()

			authHandler(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
			}

			var response transfer.APIResponse
			json.NewDecoder(w.Body).Decode(&response)

			if response.Success {
				t.Error("Expected success to be false")
			}

			if response.Message != middleware.ErrMsgUnauthorized {
				t.Errorf("Expected message '%v', got '%s'", middleware.ErrMsgUnauthorized, response.Message)
			}
		})
	}
}
