package services

import (
	"elotuschallenge/models"
	"io"
)

type IFileService interface {
	SaveFileMetadata(metadata *models.FileMetadata) (*models.FileMetadata, error)
	GetFilesByUser(userID int) ([]*models.FileMetadata, error)
	GetFileByID(fileID int) (*models.FileMetadata, error)
	SaveUploadedFile(file io.Reader, originalFilename string, contentType string, size int64, userID int, userAgent string, ipAddress string) (*models.FileMetadata, error)
}
