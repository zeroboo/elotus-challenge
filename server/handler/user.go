package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"elotuschallenge/common"
	"elotuschallenge/internal"
	"elotuschallenge/transfer"
	"elotuschallenge/utils"

	"github.com/rs/zerolog/log"
)

// handleError logs error details and response a uniform error response to client
func handleError(w http.ResponseWriter, statusCode int, userMessage string, actualError error) {
	w.Header().Set(common.HeaderContentType, common.HeaderValueContentTypeJSON)
	w.WriteHeader(statusCode)

	// Log the actual error for debugging
	refCode := utils.GenerateRandomString(8) // Generate 8-character random hex string
	if actualError != nil {
		log.Error().
			Err(actualError).
			Str("user_message", userMessage).
			Int("status_code", statusCode).
			Str("ref_code", refCode).Msg("Request failed")
	} else {
		log.Error().
			Str("user_message", userMessage).
			Int("status_code", statusCode).
			Str("ref_code", refCode).
			Msg("Request rejected")
	}

	response := transfer.NewErrorResponse(userMessage)
	response.RefCode = refCode
	json.NewEncoder(w).Encode(response)
}

// Register handler
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleError(w, http.StatusMethodNotAllowed, common.ErrMsgMethodNotAllowed, nil)
		return
	}

	// Parse request body
	var req transfer.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, http.StatusBadRequest, common.ErrMsgBadRequest, fmt.Errorf("%w: %w", common.ErrInvalidJSON, err))
		return
	}

	// Validate input
	if err := validateRegisterRequest(&req); err != nil {
		handleError(w, http.StatusBadRequest, common.ErrMsgBadRequest, fmt.Errorf("%w: %w", common.ErrInvalidRequest, err))
		return
	}

	// Check if user already exists
	exists, err := internal.UserService.UserExists(req.Username)
	if err != nil {
		handleError(w, http.StatusInternalServerError, common.ErrMsgInternalServerError, err)
		return
	}
	if exists {
		handleError(w, http.StatusConflict, common.ErrMsgUserExists, nil)
		return
	}

	// Register user (service handles password hashing)
	createdUser, err := internal.UserService.RegisterUser(req.Username, req.Password)
	if err != nil {
		handleError(w, http.StatusInternalServerError, common.ErrMsgInternalServerError, fmt.Errorf("%w: %v", common.ErrUserCreationFailed, err))
		return
	}

	// Respond with success
	w.Header().Set(common.HeaderContentType, common.HeaderValueContentTypeJSON)
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
		handleError(w, http.StatusMethodNotAllowed, common.ErrMsgMethodNotAllowed, nil)
		return
	}

	// Parse request body
	var req transfer.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, http.StatusBadRequest, common.ErrMsgBadRequest, fmt.Errorf("%w: %w", common.ErrInvalidJSON, err))
		return
	}

	// Validate input
	if err := ValidateLoginRequest(req.Username, req.Password); err != nil {
		handleError(w, http.StatusBadRequest, common.ErrMsgBadRequest, fmt.Errorf("%w: %w", common.ErrInvalidRequest, err))
		return
	}

	// Authenticate user
	user, err := internal.UserService.LoginUser(req.Username, req.Password)
	if err != nil {
		handleError(w, http.StatusUnauthorized, common.ErrMsgInvalidCredentials, fmt.Errorf("%w: %v", common.ErrInvalidCredentials, err))
		return
	}

	// Generate JWT token using standalone JWT service
	token, err := internal.TokenManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	// Respond with success
	w.Header().Set(common.HeaderContentType, common.HeaderValueContentTypeJSON)
	w.WriteHeader(http.StatusOK)

	data := transfer.LoginData{
		Auth: transfer.LoginResponse{
			Token: token,
			User: transfer.UserInfo{
				ID:       user.ID,
				Username: user.Username,
			},
		},
	}

	response := transfer.NewSuccessResponse("", data)
	json.NewEncoder(w).Encode(response)
}

// ValidateLoginRequest validates the login request
func ValidateLoginRequest(username, password string) error {
	if strings.TrimSpace(username) == "" {
		return fmt.Errorf("username is required")
	}
	if strings.TrimSpace(password) == "" {
		return fmt.Errorf("password is required")
	}
	return nil
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
