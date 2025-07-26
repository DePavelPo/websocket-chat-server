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

}
