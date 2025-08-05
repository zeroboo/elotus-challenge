package internal

import (
	"os"

	"elotuschallenge/repository"
	"elotuschallenge/services"
)

var (
	UserService  services.IUserService
	TokenManager services.ITokenManager
)

// InitServices initializes all services with their dependencies
func InitServices() {
	// Initialize repositories
	userRepo := repository.NewSQLiteUserRepository()

	// Get JWT secret from environment or use default for development
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "elotus-challenge-secret-key-change-in-production"
	}

	// Initialize services with repositories
	UserService = services.NewUserService(userRepo)
	TokenManager = services.NewTokenManager(jwtSecret, 86400) // 24 hours
}
