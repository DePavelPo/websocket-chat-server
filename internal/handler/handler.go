package handler

import (
	"github.com/DePavelPo/websocket-chat-server/internal/auth"
)

type Handler struct {
	allowedOrigins []string
	authClient     auth.Auth
}

func NewHandler(
	authCl auth.Auth,
	origins []string) *Handler {
	return &Handler{
		authClient:     authCl,
		allowedOrigins: origins,
	}
}
