package handler

import (
	"github.com/DePavelPo/websocket-chat-server/internal/auth"
)

// AuthService interface for dependency injection
type AuthService interface {
	Login(username, password string) (string, error)
	Register(username, password string) error
}

type Handler struct {
	allowedOrigins []string
	authClient     auth.Auth
	authService    AuthService
}

func NewHandler(
	authCl auth.Auth,
	authSvc AuthService,
	origins []string,
) *Handler {
	return &Handler{
		authClient:     authCl,
		authService:    authSvc,
		allowedOrigins: origins,
	}
}
