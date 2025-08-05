package models

import (
	"time"
)

// FileMetadata represents file information stored in database
type FileMetadata struct {
	ID           int       `json:"id"`
	Filename     string    `json:"filename"`
	OriginalName string    `json:"original_name"`
	ContentType  string    `json:"content_type"`
	Size         int64     `json:"size"`
	UserID       int       `json:"user_id"`
	UploadPath   string    `json:"upload_path"`
	UserAgent    string    `json:"user_agent"`
	IPAddress    string    `json:"ip_address"`
	CreatedAt    time.Time `json:"created_at"`
}

// UploadResponse represents the file upload response
type UploadResponse struct {
	Message  string       `json:"message"`
	FileInfo FileMetadata `json:"file_info"`
}
