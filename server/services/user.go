package services

import (
	"elotuschallenge/models"
	"elotuschallenge/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.IUser
}

func NewUserService(userRepo repository.IUser) IUserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// RegisterUser handles user registration with password hashing
func (s *UserService) RegisterUser(username, password string) (*models.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user model
	user := &models.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	// Save to database
	return s.userRepo.CreateUser(user)
}

// CreateUser delegates to repository (for internal use)
func (s *UserService) CreateUser(user *models.User) (*models.User, error) {
	return s.userRepo.CreateUser(user)
}

// UserExists delegates to repository
func (s *UserService) UserExists(username string) (bool, error) {
	return s.userRepo.UserExists(username)
}

// GetUserByUsername delegates to repository
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return s.userRepo.GetUserByUsername(username)
}
