package main

import (
	"log"
	"net/http"
	"os"

	"github.com/DePavelPo/websocket-chat-server/internal/auth"
	ctrlr "github.com/DePavelPo/websocket-chat-server/internal/controller"
	hl "github.com/DePavelPo/websocket-chat-server/internal/handler"
	mw "github.com/DePavelPo/websocket-chat-server/internal/middleware"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type config struct {
	Addr           string   `envconfig:"ADDR" required:"true"`
	AllowedOrigins []string `envconfig:"ALLOWED_ORIGINS" required:"true"`
	JWTKey         string   `envconfig:"JWT_KEY" required:"true"`
}

func loadConfig() (config, error) {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		logrus.WithError(err).Fatal("failed to load config")
		return config{}, err
	}
	if cfg.Addr == "" || len(cfg.AllowedOrigins) == 0 {
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

	authClient := auth.NewAuthClient(cfg.JWTKey)
	handler := hl.NewHandler(authClient, cfg.AllowedOrigins)

	_, err = os.Stat("src/index.html")
	if err != nil {
		log.Fatalf("File not found: %v", err)
	}
	fs := http.FileServer(http.Dir("src"))
	http.Handle("/chat/", http.StripPrefix("/chat/", fs))

	// Create an HTTP handler that wraps the EchoWS method
	wsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("Got WebSocket request!")
		handler.EchoWS(hub, w, r)
	})
	// Apply middleware to the WebSocket handler
	protectedWSHandler := mw.AuthMiddleware(authClient, wsHandler)

	http.Handle("/chat/ws/", protectedWSHandler)

	log.Println("Server started on", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, nil))
}
