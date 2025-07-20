package main

import (
	"log"
	"net/http"

	hl "github.com/DePavelPo/websocket-chat-server/internal/handler"
)

func main() {
	handler := hl.NewHandler()

	http.HandleFunc("/echo", handler.EchoWS)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
