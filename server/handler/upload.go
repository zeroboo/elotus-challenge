package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"elotuschallenge/common"
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
		handleError(w, http.StatusMethodNotAllowed, common.ErrMsgMethodNotAllowed, nil)
		return
	}

	// Get authenticated user from context
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		handleError(w, http.StatusUnauthorized, common.ErrMsgUserNotAuthenticated, nil)
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
		handleError(w, http.StatusBadRequest, common.ErrMsgBadRequest, err)
		return
	}

	// Get the file from form field "data"
	file, fileHeader, err := r.FormFile("data")
	if err != nil {
		handleError(w, http.StatusBadRequest, common.ErrMsgBadRequest, fmt.Errorf("failed to get file from form: %w", err))
		return
	}
	defer file.Close()

	// Check file size
	if fileHeader.Size > MaxFileSize {
		err := fmt.Errorf("%w: %v>%v", common.ErrFileTooLarge, fileHeader.Size, MaxFileSize)
		handleError(w, http.StatusBadRequest, common.ErrMsgBadRequest, err)
		return
	}

	// Detect content type
	contentType := fileHeader.Header.Get(common.HeaderContentType)
	if contentType == "" {
		// If content type is not provided, try to detect it from the file content
		buffer := make([]byte, 512)
		_, err := file.Read(buffer)
		if err != nil {
			handleError(w, http.StatusInternalServerError, common.ErrMsgReadFileFail, err)
			return
		}
		contentType = http.DetectContentType(buffer)

	}

	// Validate content type is an image
	if !IsImageContentType(contentType) {
		err := fmt.Errorf("%w: %s", common.ErrFileContentType, contentType)
		handleError(w, http.StatusBadRequest, common.ErrMsgBadRequest, err)
		return
	}

	// Get client information
	clientIP := utils.GetClientIP(r)
	userAgent := r.Header.Get(common.HeaderUserAgent)

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
		handleError(w, http.StatusInternalServerError, common.ErrMsgInternalServerError, fmt.Errorf("%w: %w", common.ErrSaveFileFail, err))
		return
	}

	log.Info().
		Int("user_id", userID).
		Str("filename", savedMetadata.Filename).
		Str("original_name", savedMetadata.OriginalName).
		Str("content_type", savedMetadata.ContentType).
		Int64("size", savedMetadata.Size).
		Str("ip_address", clientIP).
		Msg(common.MsgFileUploadSuccess)

	// Respond with success
	w.Header().Set(common.HeaderContentType, common.HeaderValueContentTypeJSON)
	w.WriteHeader(http.StatusCreated)

	data := transfer.UploadResponse{
		Message:  common.MsgFileUploadSuccess,
		FileInfo: *savedMetadata,
	}

	response := transfer.NewSuccessResponse(common.MsgFileUploadSuccess, data)
	json.NewEncoder(w).Encode(response)
}
