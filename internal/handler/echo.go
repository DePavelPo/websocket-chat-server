package handler

import (
	"net/http"
	"time"

	"github.com/DePavelPo/websocket-chat-server/internal/controller"
	"github.com/DePavelPo/websocket-chat-server/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

func newUpgrader(allowedOrigins []string) *websocket.Upgrader {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			for _, allowedOrigin := range allowedOrigins {
				if allowedOrigin == origin {
					return true
				}
			}
			return false
		},
	}
	return &upgrader
}

func (h *Handler) EchoWS(hub *controller.Hub, w http.ResponseWriter, r *http.Request) {
	upgrader := newUpgrader(h.allowedOrigins)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.WithError(err).Error("websocket upgrade error")
		return
	}

	// Wait for authentication message
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	_, message, err := conn.ReadMessage()
	if err != nil {
		logrus.WithError(err).Error("failed to read auth message")
		conn.Close()
		return
	}

	// Reset read deadline
	conn.SetReadDeadline(time.Time{})

	// Parse authentication message
	var authMsg models.AuthMessage
	if err := jsoniter.Unmarshal(message, &authMsg); err != nil {
		logrus.WithError(err).Error("failed to parse auth message")
		conn.Close()
		return
	}

	// Validate JWT token
	if authMsg.Type != "auth" || authMsg.Token == "" {
		logrus.Error("invalid auth message")
		conn.Close()
		return
	}

	username, err := h.authClient.ValidateJWT(authMsg.Token)
	if err != nil {
		logrus.WithError(err).Error("JWT validation failed")
		// Send auth failure response
		response := models.AuthResponse{Type: "auth_success", Success: false, Message: "Invalid token"}
		responseBytes, _ := jsoniter.Marshal(response)
		conn.WriteMessage(websocket.TextMessage, responseBytes)
		conn.Close()
		return
	}

	// Send auth success response
	response := models.AuthResponse{Type: "auth_success", Success: true}
	responseBytes, _ := jsoniter.Marshal(response)
	conn.WriteMessage(websocket.TextMessage, responseBytes)

	// Create client
	client := &controller.Client{
		ID:       uuid.New().String(),
		Nickname: username,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Hub:      hub,
	}

	hub.Register <- client

	go client.StartPing(30 * time.Second)
	go client.ReadProc()
	go client.WriteProc()
}
