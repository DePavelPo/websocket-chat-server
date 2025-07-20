package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) EchoWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error("websocket upgrade error: ", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logrus.Error("read message error: ", err)
			break
		}
		logrus.Infof("Received msg: %s", msg)

		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			logrus.Error("write message error: ", err)
			break
		}
	}
}
