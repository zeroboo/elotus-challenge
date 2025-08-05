package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"elotuschallenge/models"
	"elotuschallenge/repository"
	"elotuschallenge/utils"

	"github.com/rs/zerolog/log"
)

type FileService struct {
	fileRepo repository.IFile
	tmpDir   string
}

func NewFileService(fileRepo repository.IFile, tempDir string) IFileService {
	service := &FileService{
		fileRepo: fileRepo,
		tmpDir:   tempDir,
	}

	errInit := service.Init()
	if errInit != nil {
		log.Panic().Err(errInit).Msg("Failed to initialize FileService")
	}
	return service
}

func (s *FileService) Init() error {
	// Ensure the temporary directory exists
	if _, err := os.Stat(s.tmpDir); os.IsNotExist(err) {
		err = os.MkdirAll(s.tmpDir, 0755)
		if err != nil {
			log.Error().Err(err).Str("path", s.tmpDir).Msg("Failed to create temporary directory")
			return fmt.Errorf("failed to create temporary directory: %w", err)
		}
		log.Info().Str("path", s.tmpDir).Msg("Temporary directory created")
	} else {
		log.Info().Str("path", s.tmpDir).Msg("Temporary directory already exists")
	}

	return nil
}

// SaveFileMetadata saves file metadata to database
func (s *FileService) SaveFileMetadata(metadata *models.FileMetadata) (*models.FileMetadata, error) {
	return s.fileRepo.CreateFile(metadata)
}

// GetFilesByUser retrieves all files for a specific user
func (s *FileService) GetFilesByUser(userID int) ([]*models.FileMetadata, error) {
	return s.fileRepo.GetFilesByUser(userID)
}

// GetFileByID retrieves a specific file by ID
func (s *FileService) GetFileByID(fileID int) (*models.FileMetadata, error) {
	return s.fileRepo.GetFileByID(fileID)
}

// SaveUploadedFile handles the complete process of saving a file to disk and database
func (s *FileService) SaveUploadedFile(file io.Reader, originalFilename string, contentType string, size int64, userID int, userAgent string, ipAddress string) (*models.FileMetadata, error) {
	// Generate unique filename
	uniqueFilename := fmt.Sprintf("%s_%s%s",
		utils.GenerateRandomString(12),
		time.Now().Format("20060102_150405"),
		filepath.Ext(originalFilename))

	// Create full path in /tmp directory
	tmpFilePath := filepath.Join(s.tmpDir, uniqueFilename)

	// Create the temporary file
	tmpFile, err := os.Create(tmpFilePath)
	if err != nil {
		log.Error().Err(err).Str("path", tmpFilePath).Msg("Failed to create temporary file")
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tmpFile.Close()

	// Copy uploaded file content to temporary file
	bytesWritten, err := io.Copy(tmpFile, file)
	if err != nil {
		log.Error().Err(err).Msg("Failed to write file content")
		// Clean up the temporary file
		os.Remove(tmpFilePath)
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Verify the written size matches the uploaded size
	if bytesWritten != size {
		log.Error().Int64("written", bytesWritten).Int64("expected", size).Msg("File size mismatch")
		os.Remove(tmpFilePath)
		return nil, fmt.Errorf("file upload incomplete: expected %d bytes, wrote %d bytes", size, bytesWritten)
	}

	// Create file metadata
	fileMetadata := &models.FileMetadata{
		Filename:     uniqueFilename,
		OriginalName: originalFilename,
		ContentType:  contentType,
		Size:         size,
		UserID:       userID,
		UploadPath:   tmpFilePath,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		CreatedAt:    time.Now(),
	}

	// Save metadata to database
	savedMetadata, err := s.fileRepo.CreateFile(fileMetadata)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save file metadata to database")
		// Clean up the temporary file if database save fails
		errRemoveFile := os.Remove(tmpFilePath)
		if errRemoveFile != nil {
			log.Error().Err(errRemoveFile).Msg("Failed to remove temporary file after metadata save failure")
		}
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}

	log.Info().
		Int("user_id", userID).
		Str("filename", uniqueFilename).
		Str("original_name", originalFilename).
		Str("content_type", contentType).
		Int64("size", size).
		Str("ip_address", ipAddress).
		Msg("Success")

	return savedMetadata, nil
}
