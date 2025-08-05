package internal

import (
	"elotuschallenge/repository"
	"elotuschallenge/services"
)

var (
	UserService services.IUserService
)

// InitServices initializes all services with their dependencies
func InitServices() {
	// Initialize repositories
	userRepo := repository.NewSQLiteUserRepository()

	// Initialize services with repositories
	UserService = services.NewUserService(userRepo)
}
