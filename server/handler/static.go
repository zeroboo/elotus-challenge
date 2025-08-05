package handler

import (
	"elotuschallenge/common"
	"net/http"
	"path/filepath"
)

// HandleStatic serves static files
func HandleStatic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, http.StatusMethodNotAllowed, common.ErrMsgMethodNotAllowed, nil)
		return
	}

	// Serve the upload form
	if r.URL.Path == "/upload-form" || r.URL.Path == "/upload-form/" {
		http.ServeFile(w, r, filepath.Join("static", "upload.html"))
		return
	}

	http.NotFound(w, r)
}
