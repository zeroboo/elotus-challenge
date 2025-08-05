package repository

import "elotuschallenge/models"

type IFile interface {
	CreateFile(file *models.FileMetadata) (*models.FileMetadata, error)
	GetFileByID(fileID int) (*models.FileMetadata, error)
	GetFilesByUser(userID int) ([]*models.FileMetadata, error)
}
