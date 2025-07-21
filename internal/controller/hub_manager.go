package controller

import "github.com/DePavelPo/websocket-chat-server/internal/models"

type Hub struct {
	Clients    map[*models.Client]bool
	Register   chan *models.Client
	Unregister chan *models.Client
	Broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*models.Client]bool),
		Register:   make(chan *models.Client),
		Unregister: make(chan *models.Client),
		Broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					delete(h.Clients, client)
					close(client.Send)
				}
			}
		}
	}
}
