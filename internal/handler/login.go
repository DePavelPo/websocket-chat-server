package handler

import (
	"net/http"

	"github.com/DePavelPo/websocket-chat-server/internal/models"
	"github.com/sirupsen/logrus"

	jsoniter "github.com/json-iterator/go"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Login request received")

	var creds models.Credentials
	if err := jsoniter.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if creds.Username == "" || creds.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Use the auth service for login
	token, err := h.authService.Login(creds.Username, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Return the token
	w.Header().Set("Content-Type", "application/json")
	jsoniter.NewEncoder(w).Encode(models.LoginResponse{Token: token})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Register request received")

	var creds models.Credentials
	if err := jsoniter.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if creds.Username == "" || creds.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Use the auth service for registration
	err := h.authService.Register(creds.Username, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}
