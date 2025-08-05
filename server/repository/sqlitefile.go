package repository

import (
	"database/sql"
	"elotuschallenge/database"
	"elotuschallenge/models"
)

type SQLiteFileRepository struct{}

func NewSQLiteFileRepository() IFile {
	return &SQLiteFileRepository{}
}

// CreateFile inserts a new file metadata into the database
func (r *SQLiteFileRepository) CreateFile(file *models.FileMetadata) (*models.FileMetadata, error) {
	query := `
		INSERT INTO files (filename, original_name, content_type, size, user_id, upload_path, user_agent, ip_address, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`

	result, err := database.DB.Exec(query, file.Filename, file.OriginalName, file.ContentType, file.Size, file.UserID, file.UploadPath, file.UserAgent, file.IPAddress)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	file.ID = int(id)
	return file, nil
}

// GetFileByID retrieves a file by its ID
func (r *SQLiteFileRepository) GetFileByID(fileID int) (*models.FileMetadata, error) {
	query := "SELECT id, filename, original_name, content_type, size, user_id, upload_path, user_agent, ip_address, created_at FROM files WHERE id = ?"
	var file models.FileMetadata
	err := database.DB.QueryRow(query, fileID).Scan(&file.ID, &file.Filename, &file.OriginalName, &file.ContentType, &file.Size, &file.UserID, &file.UploadPath, &file.UserAgent, &file.IPAddress, &file.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // File not found
		}
		return nil, err
	}
	return &file, nil
}

// GetFilesByUser retrieves all files for a specific user
func (r *SQLiteFileRepository) GetFilesByUser(userID int) ([]*models.FileMetadata, error) {
	query := "SELECT id, filename, original_name, content_type, size, user_id, upload_path, user_agent, ip_address, created_at FROM files WHERE user_id = ? ORDER BY created_at DESC"
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*models.FileMetadata
	for rows.Next() {
		var file models.FileMetadata
		err := rows.Scan(&file.ID, &file.Filename, &file.OriginalName, &file.ContentType, &file.Size, &file.UserID, &file.UploadPath, &file.UserAgent, &file.IPAddress, &file.CreatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}

	return files, nil
}
