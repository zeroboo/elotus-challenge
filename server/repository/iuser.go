package repository

import "elotuschallenge/models"

type IUser interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UserExists(username string) (bool, error)
}
