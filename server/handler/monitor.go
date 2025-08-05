package handler

import (
	"elotuschallenge/common"
	"fmt"
	"net/http"
)

// Health check handler
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, http.StatusMethodNotAllowed, common.ErrMsgMethodNotAllowed, nil)
		return
	}

	w.Header().Set(common.HeaderContentType, common.HeaderValueContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"ok","message":"Server is running"}`)
}
