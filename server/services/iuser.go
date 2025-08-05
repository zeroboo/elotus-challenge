package services

import "elotuschallenge/models"

type IUserService interface {
	CreateUser(user *models.User) (*models.User, error)
	UserExists(username string) (bool, error)
	GetUserByUsername(username string) (*models.User, error)
}
