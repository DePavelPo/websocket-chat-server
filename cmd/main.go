package main

import (
	"log"
	"net/http"
	"os"

	ctrlr "github.com/DePavelPo/websocket-chat-server/internal/controller"
	hl "github.com/DePavelPo/websocket-chat-server/internal/handler"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	host := getEnvOrDefault("ADDR", "localhost:8080")

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

	log.Println("Server started on", host)
	log.Fatal(http.ListenAndServe(host, nil))
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found or failed to load")
	}
}

func getEnvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
