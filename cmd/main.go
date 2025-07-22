package main

import (
	"log"
	"net/http"
	"os"

	ctrlr "github.com/DePavelPo/websocket-chat-server/internal/controller"
	hl "github.com/DePavelPo/websocket-chat-server/internal/handler"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type config struct {
	Addr string `envconfig:"ADDR" required:"true"`
}

func loadConfig() (config, error) {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		logrus.WithError(err).Fatal("failed to load config")
		return config{}, err
	}
	if cfg.Addr == "" {
		return config{}, envconfig.ErrInvalidSpecification
	}
	return cfg, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		logrus.Fatalf("missing or invalid config: %v", err)
	}

	hub := ctrlr.NewHub()
	go hub.Run()

	handler := hl.NewHandler()

	_, err = os.Stat("src/index.html")
	if err != nil {
		log.Fatalf("File not found: %v", err)
	}
	fs := http.FileServer(http.Dir("src"))
	http.Handle("/chat/", http.StripPrefix("/chat/", fs))

	http.HandleFunc("/chat/ws", func(w http.ResponseWriter, r *http.Request) {
		handler.EchoWS(hub, w, r)
	})

	log.Println("Server started on", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, nil))
}
