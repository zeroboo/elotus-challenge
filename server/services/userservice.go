package services

import (
	"elotuschallenge/models"
	"elotuschallenge/repository"
	"fmt"

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

// LoginUser authenticates a user with username and password
func (s *UserService) LoginUser(username, password string) (*models.User, error) {
	// Get user by username
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
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
