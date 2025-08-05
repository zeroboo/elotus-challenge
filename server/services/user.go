package services

import (
	"elotuschallenge/models"
	"elotuschallenge/repository"
)

type UserService struct {
	userRepo repository.IUser
}

func NewUserService(userRepo repository.IUser) IUserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser delegates to repository
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
