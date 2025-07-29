package main

import (
	"log"
	"net/http"
	"os"

	"github.com/DePavelPo/websocket-chat-server/db/repository"
	"github.com/DePavelPo/websocket-chat-server/internal/auth"
	ctrlr "github.com/DePavelPo/websocket-chat-server/internal/controller"
	hl "github.com/DePavelPo/websocket-chat-server/internal/handler"
	"github.com/DePavelPo/websocket-chat-server/internal/service"
	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/sirupsen/logrus"
)

type config struct {
	Addr           string   `envconfig:"ADDR" required:"true"`
	AllowedOrigins []string `envconfig:"ALLOWED_ORIGINS" required:"true"`
	JWTKey         string   `envconfig:"JWT_KEY" required:"true"`
	DBConn         string   `envconfig:"DB_CONN" required:"true"`
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

	// Connect to database
	db, err := initDB(cfg.DBConn)
	if err != nil {
		logrus.WithError(err).Fatal("failed to initialize database")
	}
	defer db.Close()

	// Initialize repositories
	repo := repository.NewRepository(db)
	userRepo := repo.GetUserRepository()
	sessionRepo := repo.GetSessionRepository()

	// Initialize auth client and service
	authClient := auth.NewAuthClient(cfg.JWTKey)
	authService := service.NewAuthService(userRepo, sessionRepo, authClient)
	handler := hl.NewHandler(authClient, authService, cfg.AllowedOrigins)

	_, err = os.Stat("src/index.html")
	if err != nil {
		log.Fatalf("File not found: %v", err)
	}
	fs := http.FileServer(http.Dir("src"))
	http.Handle("/chat/", http.StripPrefix("/chat/", fs))

	// Add login route
	http.HandleFunc("/chat/login", handler.Login)
	http.HandleFunc("/chat/register", handler.Register)

	// WebSocket route - authentication handled in the handler
	http.HandleFunc("/chat/ws/", func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("Got WebSocket request!")
		handler.EchoWS(hub, w, r)
	})

	log.Println("Server started on", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, nil))
}

func initDB(dbConn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dbConn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, errors.New("cannot ping db")
	}
	return db, nil
}
