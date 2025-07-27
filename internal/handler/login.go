package handler

import (
	"net/http"

	"github.com/DePavelPo/websocket-chat-server/internal/models"

	jsoniter "github.com/json-iterator/go"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	if err := jsoniter.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if creds.Username == "" || creds.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Generate JWT token
	token, err := h.authClient.GenerateJWT(creds.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return the token
	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(models.LoginResponse{Token: token})
}
