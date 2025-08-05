package handler

import (
	"encoding/json"
	"net/http"

	"elotuschallenge/transfer"

	"github.com/rs/zerolog/log"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := transfer.NewErrorResponse("Method not allowed")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get authenticated user from context
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		response := transfer.NewErrorResponse("User not authenticated")
		json.NewEncoder(w).Encode(response)
		return
	}

	username, _ := r.Context().Value("username").(string)

	log.Info().
		Int("user_id", userID).
		Str("username", username).
		Msg("Authenticated user accessed upload endpoint")

	// TODO: Implement file upload logic
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	data := map[string]interface{}{
		"message": "Upload endpoint - authentication working",
		"user":    map[string]interface{}{"id": userID, "username": username},
	}

	response := transfer.NewSuccessResponse("Upload endpoint accessible", data)
	json.NewEncoder(w).Encode(response)
}
