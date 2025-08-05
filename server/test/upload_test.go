package test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"

	"elotuschallenge/common"
	"elotuschallenge/handler"
	"elotuschallenge/middleware"
	"elotuschallenge/test/share"
	"elotuschallenge/transfer"
)

func TestHandleUpload_ValidImageFile_Success(t *testing.T) {
	// First register and login to get a token
	token := loginTestUser(t, "uploaduser", "password123")

	// Load PNG data from file
	pngData, err := share.LoadTestPNG("./test/files/leaf.png")
	if err != nil {
		t.Fatalf("Failed to load test PNG file: %v", err)
	}

	// Create multipart form
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add file field with explicit Content-Type header
	partHeader := textproto.MIMEHeader{}
	partHeader.Set("Content-Disposition", `form-data; name="data"; filename="leaf.png"`)
	partHeader.Set(common.HeaderContentType, "image/png")
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		t.Fatalf("Failed to create form part: %v", err)
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
	req.Header.Set(common.HeaderContentType, writer.FormDataContentType())
	req.Header.Set(common.HeaderAuthorization, "Bearer "+token)

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
	req.Header.Set(common.HeaderContentType, writer.FormDataContentType())
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
	req.Header.Set(common.HeaderContentType, writer.FormDataContentType())
	req.Header.Set(common.HeaderAuthorization, "Bearer "+token)

	w := httptest.NewRecorder()

	authHandler := middleware.AuthUser(handler.HandleUpload)
	authHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response transfer.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Message != common.ErrMsgBadRequest {
		t.Errorf("Expected '%v', got '%s'", common.ErrMsgBadRequest, response.Message)
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
	req.Header.Set(common.HeaderContentType, writer.FormDataContentType())
	req.Header.Set(common.HeaderAuthorization, "Bearer "+token)

	w := httptest.NewRecorder()

	authHandler := middleware.AuthUser(handler.HandleUpload)
	authHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response transfer.APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Message != common.ErrMsgBadRequest {
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
	req.Header.Set(common.HeaderContentType, common.HeaderValueContentTypeJSON)
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
