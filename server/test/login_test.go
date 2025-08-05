package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"elotuschallenge/common"
	"elotuschallenge/handler"
	"elotuschallenge/transfer"
)

const ContentTypeJSON = "application/json"
const HeaderContentType = "Content-Type"

func TestHandleLogin_ValidCredentials_Success(t *testing.T) {
	// First register a user
	registerUser(t, "loginuser", "password123")

	// Now login with the same credentials
	reqBody := transfer.LoginRequest{
		Username: "loginuser",
		Password: "password123",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set(HeaderContentType, ContentTypeJSON)
	w := httptest.NewRecorder()

	handler.HandleLogin(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response
	var response transfer.APIResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Check response structure
	if !response.Success {
		t.Error("Expected success to be true")
	}

	if response.Data == nil {
		t.Fatal("Expected data to be present")
	}

	// Check login data structure
	dataMap, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Expected data to be a map")
	}

	auth, ok := dataMap["auth"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected auth to be present in data")
	}

	// Check token is present
	if token, ok := auth["token"].(string); !ok || token == "" {
		t.Error("Expected non-empty token in response")
	}

	// Check user info
	user, ok := auth["user"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected user to be present in auth")
	}

	if username, ok := user["username"].(string); !ok || username != "loginuser" {
		t.Errorf("Expected username 'loginuser', got '%v'", user["username"])
	}
}

func TestHandleLogin_InvalidCredentials_Error(t *testing.T) {
	// First register a user
	registerUser(t, "loginuser2", "password123")

	// Try to login with wrong password
	reqBody := transfer.LoginRequest{
		Username: "loginuser2",
		Password: "wrongpassword",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set(HeaderContentType, ContentTypeJSON)
	w := httptest.NewRecorder()

	handler.HandleLogin(w, req)

	// Check status code
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	// Parse response
	var response transfer.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Message != "Invalid credentials" {
		t.Errorf("Expected message 'Invalid credentials', got '%s'", response.Message)
	}
}

func TestHandleLogin_NonexistentUser_Error(t *testing.T) {
	reqBody := transfer.LoginRequest{
		Username: "nonexistentuser",
		Password: "password123",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set(HeaderContentType, ContentTypeJSON)
	w := httptest.NewRecorder()

	handler.HandleLogin(w, req)

	// Check status code
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	// Parse response
	var response transfer.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Message != "Invalid credentials" {
		t.Errorf("Expected message 'Invalid credentials', got '%s'", response.Message)
	}
}

func TestHandleLogin_InvalidMethod_Error(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	handler.HandleLogin(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}

	var response transfer.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Message != common.ErrMsgMethodNotAllowed {
		t.Errorf("Expected message '%v', got '%s'", common.ErrMsgMethodNotAllowed, response.Message)
	}
}

func TestHandleLogin_InvalidJSON_Error(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("invalid json")))
	req.Header.Set(HeaderContentType, ContentTypeJSON)
	w := httptest.NewRecorder()

	handler.HandleLogin(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response transfer.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Message != common.ErrMsgBadRequest {
		t.Errorf("Expected message '%v', got '%s'", common.ErrMsgBadRequest, response.Message)
	}
}

func TestHandleLogin_ValidationErrors_Error(t *testing.T) {
	testCases := []struct {
		name     string
		username string
		password string
		expected string
	}{
		{
			name:     "Empty username",
			username: "",
			password: "password123",
			expected: common.ErrMsgBadRequest,
		},
		{
			name:     "Empty password",
			username: "testuser",
			password: "",
			expected: common.ErrMsgBadRequest,
		},
		{
			name:     "Both empty",
			username: "",
			password: "",
			expected: common.ErrMsgBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody := transfer.LoginRequest{
				Username: tc.username,
				Password: tc.password,
			}

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			req.Header.Set(HeaderContentType, ContentTypeJSON)
			w := httptest.NewRecorder()

			handler.HandleLogin(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
			}

			var response transfer.APIResponse
			json.NewDecoder(w.Body).Decode(&response)

			if response.Success {
				t.Error("Expected success to be false")
			}

			if response.Message != tc.expected {
				t.Errorf("Expected message '%s', got '%s'", tc.expected, response.Message)
			}
		})
	}
}

// Helper function to register a user for testing
func registerUser(t *testing.T, username, password string) {
	reqBody := transfer.RegisterRequest{
		Username: username,
		Password: password,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set(HeaderContentType, ContentTypeJSON)
	w := httptest.NewRecorder()

	handler.HandleRegister(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to register user for test: %s", w.Body.String())
	}
}
