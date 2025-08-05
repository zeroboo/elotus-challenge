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
	if r.URL.Path == "/form/upload" || r.URL.Path == "/form/upload/" {
		http.ServeFile(w, r, filepath.Join("static", "upload.html"))
		return
	}

	// Serve the registration form
	if r.URL.Path == "/form/register" || r.URL.Path == "/form/register/" {
		http.ServeFile(w, r, filepath.Join("static", "register.html"))
		return
	}

	// Serve the login form
	if r.URL.Path == "/form/login" || r.URL.Path == "/form/login/" {
		http.ServeFile(w, r, filepath.Join("static", "login.html"))
		return
	}

	http.NotFound(w, r)
}
