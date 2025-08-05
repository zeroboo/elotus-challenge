package internal

import (
	"os"
	"strconv"

	"elotuschallenge/repository"
	"elotuschallenge/services"
)

var (
	UserService  services.IUserService
	TokenManager services.ITokenManager
	FileService  services.IFileService
)

// InitServices initializes all services with their dependencies
func InitServices() {
	// Initialize repositories
	userRepo := repository.NewSQLiteUserRepository()
	fileRepo := repository.NewSQLiteFileRepository()

	// Get JWT secret from environment or use default for development
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "elotus-challenge-default"
	}

	// Get temp directory from environment or use default
	tempDir := os.Getenv("TEMP_DIR")
	if tempDir == "" {
		tempDir = "./tmp"
	}

	// Get token expiration from environment or use default (24 hours)
	tokenExpirationSeconds := int64(86400) // Default: 24 hours
	if tokenExpEnv := os.Getenv("TOKEN_EXPIRATION_SECONDS"); tokenExpEnv != "" {
		if expSeconds, err := strconv.ParseInt(tokenExpEnv, 10, 64); err == nil && expSeconds > 0 {
			tokenExpirationSeconds = expSeconds
		}
	}

	// Initialize services with repositories
	UserService = services.NewUserService(userRepo)
	TokenManager = services.NewTokenManager(jwtSecret, tokenExpirationSeconds)
	FileService = services.NewFileService(fileRepo, tempDir)
}
