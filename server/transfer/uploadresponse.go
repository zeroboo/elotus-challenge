package transfer

import "elotuschallenge/models"

// UploadResponse represents the file upload response
type UploadResponse struct {
	Message  string              `json:"message"`
	FileInfo models.FileMetadata `json:"file_info"`
}
