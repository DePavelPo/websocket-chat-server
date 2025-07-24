package handler

import (
	"net/http"
	"time"

	"github.com/DePavelPo/websocket-chat-server/internal/controller"
	"github.com/DePavelPo/websocket-chat-server/utils"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

	client := &controller.Client{
		ID:       uuid.New().String(),
		Nickname: utils.GenerateNickname(),
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Hub:      hub,
	}

	hub.Register <- client

	go client.StartPing(30 * time.Second)
	go client.ReadProc()
	go client.WriteProc()
}
