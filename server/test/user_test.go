package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"elotuschallenge/handler"
	"elotuschallenge/transfer"
)

func TestHandleRegister_ValidRequest_Success(t *testing.T) {
	// Create request body
	reqBody := transfer.RegisterRequest{
		Username: "testuser123",
		Password: "password123",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler.HandleRegister(w, req)

	// Check status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
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
		t.Error("Expected data to be present")
	}
}

func TestHandleRegister_InvalidMethod_ResponseError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/register", nil)
	w := httptest.NewRecorder()

	handler.HandleRegister(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}

	var response transfer.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Message != "Method not allowed" {
		t.Errorf("Expected message 'Method not allowed', got '%s'", response.Message)
	}
}

func TestHandleRegister_InvalidJSON_ResponseError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleRegister(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response transfer.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Message != "Invalid JSON format" {
		t.Errorf("Expected message 'Invalid JSON format', got '%s'", response.Message)
	}
}

func TestHandleRegister_InvalidRequest_ResponseError(t *testing.T) {
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
			expected: "Validation failed",
		},
		{
			name:     "Short username",
			username: "ab",
			password: "password123",
			expected: "Validation failed",
		},
		{
			name:     "Long username",
			username: "this_is_a_very_long_username_that_exceeds_fifty_characters_limit",
			password: "password123",
			expected: "Validation failed",
		},
		{
			name:     "Empty password",
			username: "testuser",
			password: "",
			expected: "Validation failed",
		},
		{
			name:     "Short password",
			username: "testuser",
			password: "12345",
			expected: "Validation failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody := transfer.RegisterRequest{
				Username: tc.username,
				Password: tc.password,
			}

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandleRegister(w, req)

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

func TestHandleRegister_UserAlreadyExists_ResponseError(t *testing.T) {
	// First, create a user
	reqBody := transfer.RegisterRequest{
		Username: "duplicateuser",
		Password: "password123",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Register first user
	handler.HandleRegister(w, req)

	// Try to register the same user again
	body, _ = json.Marshal(reqBody)
	req = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler.HandleRegister(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, w.Code)
	}

	var response transfer.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Message != "Username already exists" {
		t.Errorf("Expected message 'Username already exists', got '%s'", response.Message)
	}
}

func TestHandleRegister_InvalidContentType_ResponseError(t *testing.T) {
	reqBody := transfer.RegisterRequest{
		Username: "testuser456",
		Password: "password123",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleRegister(w, req)

	// Check response content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestHandleRegister_ResponseStructure_Correct(t *testing.T) {
	reqBody := transfer.RegisterRequest{
		Username: "structuretest",
		Password: "password123",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleRegister(w, req)

	var response transfer.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	// Check that data contains user information
	if response.Data == nil {
		t.Fatal("Expected data to be present")
	}

	// Convert data to map for easier testing
	dataMap, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Expected data to be a map")
	}

	user, ok := dataMap["user"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected user to be present in data")
	}

	// Check user fields
	if _, ok := user["id"]; !ok {
		t.Error("Expected user id to be present")
	}

	if username, ok := user["username"].(string); !ok || username != "structuretest" {
		t.Errorf("Expected username 'structuretest', got '%v'", user["username"])
	}

	// Ensure password is not included in response
	if _, ok := user["password"]; ok {
		t.Error("Password should not be included in response")
	}
}
