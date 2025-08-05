package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"elotuschallenge/internal"
	"elotuschallenge/transfer"
	"elotuschallenge/utils"

	"github.com/rs/zerolog/log"
)

const (
	MaxFileSize = 8 << 20 // 8 MB in bytes
)

// IsImageContentType checks if the content type is a valid image type
func IsImageContentType(contentType string) bool {
	validImageTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/bmp",
		"image/tiff",
		"image/svg+xml",
	}

	contentType = strings.ToLower(contentType)
	for _, validType := range validImageTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	// Get authenticated user from context
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		handleError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	username, _ := r.Context().Value("username").(string)

	log.Info().
		Int("user_id", userID).
		Str("username", username).
		Msg("File upload request received")

	// Parse multipart form with size limit
	err := r.ParseMultipartForm(MaxFileSize)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse multipart form")
		handleError(w, http.StatusBadRequest, "Failed to parse form data or file too large", err)
		return
	}

	// Get the file from form field "data"
	file, fileHeader, err := r.FormFile("data")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get file from form")
		handleError(w, http.StatusBadRequest, "File upload field 'data' is required", err)
		return
	}
	defer file.Close()

	// Check file size
	if fileHeader.Size > MaxFileSize {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("File size (%d bytes) exceeds maximum allowed size (%d bytes)", fileHeader.Size, MaxFileSize), nil)
		return
	}

	// Detect content type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		// If content type is not provided, try to detect it from the file content
		buffer := make([]byte, 512)
		_, err := file.Read(buffer)
		if err != nil {
			handleError(w, http.StatusInternalServerError, "Failed to read file content", err)
			return
		}
		contentType = http.DetectContentType(buffer)

		// Reset file position to beginning
		file.Seek(0, 0)
	}

	// Validate content type is an image
	if !IsImageContentType(contentType) {
		handleError(w, http.StatusBadRequest, fmt.Sprintf("Invalid file type: %s. Only image files are allowed", contentType), nil)
		return
	}

	// Get client information
	clientIP := utils.GetClientIP(r)
	userAgent := r.Header.Get("User-Agent")

	// Use file service to save the uploaded file
	savedMetadata, err := internal.FileService.SaveUploadedFile(
		file,
		fileHeader.Filename,
		contentType,
		fileHeader.Size,
		userID,
		userAgent,
		clientIP,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save uploaded file")
		handleError(w, http.StatusInternalServerError, "Failed to save file", err)
		return
	}

	log.Info().
		Int("user_id", userID).
		Str("filename", savedMetadata.Filename).
		Str("original_name", savedMetadata.OriginalName).
		Str("content_type", savedMetadata.ContentType).
		Int64("size", savedMetadata.Size).
		Str("ip_address", clientIP).
		Msg("File uploaded successfully")

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	data := transfer.UploadResponse{
		Message:  "File uploaded successfully",
		FileInfo: *savedMetadata,
	}

	response := transfer.NewSuccessResponse("File uploaded successfully", data)
	json.NewEncoder(w).Encode(response)
}
