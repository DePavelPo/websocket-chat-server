package handler

import (
	"net/http"

	"github.com/DePavelPo/websocket-chat-server/internal/controller"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) EchoWS(hub *controller.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error("websocket upgrade error: ", err)
		return
	}

	client := &controller.Client{
		Conn: conn,
		Send: make(chan []byte),
		Hub:  hub,
	}

	hub.Register <- client

	go client.ReadProc()
	go client.WriteProc()
}
