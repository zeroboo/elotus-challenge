package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"elotuschallenge/internal"
	"elotuschallenge/transfer"

	"github.com/rs/zerolog/log"
)

// handleError logs error details and response a uniform error response to client
func handleError(w http.ResponseWriter, statusCode int, userMessage string, actualError error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Log the actual error for debugging
	if actualError != nil {
		log.Error().Err(actualError).Str("user_message", userMessage).Int("status_code", statusCode).Msg("Request failed")
	} else {
		log.Warn().Str("user_message", userMessage).Int("status_code", statusCode).Msg("Request rejected")
	}

	response := transfer.NewErrorResponse(userMessage)
	json.NewEncoder(w).Encode(response)
}

// Register handler
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	// Parse request body
	var req transfer.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid JSON format", err)
		return
	}

	// Validate input
	if err := validateRegisterRequest(&req); err != nil {
		handleError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Check if user already exists
	exists, err := internal.UserService.UserExists(req.Username)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to check user existence", err)
		return
	}
	if exists {
		handleError(w, http.StatusConflict, "Username already exists", nil)
		return
	}

	// Register user (service handles password hashing)
	createdUser, err := internal.UserService.RegisterUser(req.Username, req.Password)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	data := transfer.RegisterData{
		User: transfer.UserInfo{
			ID:       createdUser.ID,
			Username: createdUser.Username,
		},
	}

	response := transfer.NewSuccessResponse("User registered successfully", data)
	json.NewEncoder(w).Encode(response)
}

// Login handler
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	// TODO: Implement login logic
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := transfer.NewSuccessResponse("Login endpoint - TODO", nil)
	json.NewEncoder(w).Encode(response)
}

// validateRegisterRequest validates the registration request
func validateRegisterRequest(req *transfer.RegisterRequest) error {
	if strings.TrimSpace(req.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if len(req.Username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}
	if len(req.Username) > 50 {
		return fmt.Errorf("username must be less than 50 characters")
	}
	if strings.TrimSpace(req.Password) == "" {
		return fmt.Errorf("password is required")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}
	return nil
}
