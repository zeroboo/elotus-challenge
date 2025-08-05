package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"elotuschallenge/internal"
	"elotuschallenge/models"
	"elotuschallenge/transfer"

	"golang.org/x/crypto/bcrypt"
)

// Register handler
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req transfer.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate input
	if err := validateRegisterRequest(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user already exists
	exists, err := internal.UserService.UserExists(req.Username)
	if err != nil {
		http.Error(w, "Failed to check user existence", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	// Create user
	user := models.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	}

	// Save to database using service
	createdUser, err := internal.UserService.CreateUser(&user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"message": "User registered successfully",
		"user": map[string]interface{}{
			"id":       createdUser.ID,
			"username": createdUser.Username,
		},
	}

	json.NewEncoder(w).Encode(response)
}

// Login handler
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Implement login logic
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"message":"Login endpoint - TODO"}`)
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
