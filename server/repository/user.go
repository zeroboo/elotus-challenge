package repository

import "elotuschallenge/models"

type IUserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByName(name string) (*models.User, error)
}
