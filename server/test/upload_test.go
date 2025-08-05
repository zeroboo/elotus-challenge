package test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"elotuschallenge/handler"
	"elotuschallenge/middleware"
	"elotuschallenge/transfer"
)

func TestHandleUpload_ValidImageFile_Success(t *testing.T) {
	// First register and login to get a token
	token := loginTestUser(t, "uploaduser", "password123")

	// Create a simple 1x1 PNG image in memory
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, // 1x1 dimensions
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53, // bit depth, color type, etc.
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41, // IDAT chunk
		0x54, 0x08, 0xD7, 0x63, 0xF8, 0x00, 0x00, 0x00, // minimal image data
		0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x37, 0x6E, 0xF9, 0x24, 0x00, 0x00, 0x00, 0x00, // IEND chunk
		0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}

	// Create multipart form
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add file field
	part, err := writer.CreateFormFile("data", "test.png")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}

	_, err = io.Copy(part, bytes.NewReader(pngData))
	if err != nil {
		t.Fatalf("Failed to write file data: %v", err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatalf("Failed to close multipart writer: %v", err)
	}

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	// Call handler through middleware
	authHandler := middleware.AuthUser(handler.HandleUpload)
	authHandler(w, req)

	// Check response
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var response transfer.APIResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}

	if !strings.Contains(response.Message, "uploaded successfully") {
		t.Errorf("Expected success message, got '%s'", response.Message)
	}
}

func TestHandleUpload_NoAuthToken_Error(t *testing.T) {
	// Create simple multipart form
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, _ := writer.CreateFormFile("data", "test.txt")
	part.Write([]byte("test content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// No Authorization header

	w := httptest.NewRecorder()

	authHandler := middleware.AuthUser(handler.HandleUpload)
	authHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestHandleUpload_InvalidFileType_Error(t *testing.T) {
	token := loginTestUser(t, "uploaduser2", "password123")

	// Create text file instead of image
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, _ := writer.CreateFormFile("data", "test.txt")
	part.Write([]byte("This is not an image"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	authHandler := middleware.AuthUser(handler.HandleUpload)
	authHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response transfer.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if !strings.Contains(response.Message, "Invalid file type") {
		t.Errorf("Expected invalid file type error, got '%s'", response.Message)
	}
}

func TestHandleUpload_NoFileField_Error(t *testing.T) {
	token := loginTestUser(t, "uploaduser3", "password123")

	// Create form without the required "data" field
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, _ := writer.CreateFormFile("wrongfield", "test.png")
	part.Write([]byte("fake image data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	authHandler := middleware.AuthUser(handler.HandleUpload)
	authHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response transfer.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if !strings.Contains(response.Message, "field 'data' is required") {
		t.Errorf("Expected missing field error, got '%s'", response.Message)
	}
}

// Helper function to register and login a test user, returning JWT token
func loginTestUser(t *testing.T, username, password string) string {
	// Register user first
	registerUser(t, username, password)

	// Login to get token
	reqBody := transfer.LoginRequest{
		Username: username,
		Password: password,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleLogin(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Failed to login test user: %s", w.Body.String())
	}

	var response transfer.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	dataMap := response.Data.(map[string]interface{})
	authMap := dataMap["auth"].(map[string]interface{})
	return authMap["token"].(string)
}
