package services

import (
	"elotuschallenge/models"
)

type IUserService interface {
	RegisterUser(username, password string) (*models.User, error)
	LoginUser(username, password string) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	UserExists(username string) (bool, error)
	GetUserByUsername(username string) (*models.User, error)
}
