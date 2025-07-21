package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	ctrlr "github.com/DePavelPo/websocket-chat-server/internal/controller"
	hl "github.com/DePavelPo/websocket-chat-server/internal/handler"
)

func main() {
	hub := ctrlr.NewHub()
	go hub.Run()

	handler := hl.NewHandler()

	_, err := os.Stat("src/index.html")
	if err != nil {
		log.Fatalf("File not found: %v", err)
	}
	fs := http.FileServer(http.Dir("src"))
	http.Handle("/", fs)

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		handler.EchoWS(hub, w, r)
	})
	port := 8080
	log.Printf("Server started on :%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
