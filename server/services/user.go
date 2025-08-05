package services

import "elotuschallenge/models"

type IUserService interface {
	CreateUser(user *models.User) (*models.User, error)
}
