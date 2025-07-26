package middleware

import (
	"context"
	"net/http"

	"github.com/DePavelPo/websocket-chat-server/internal/auth"
)

type ctxKey string

var ctxUsernameKey ctxKey = "username"

func AuthMiddleware(auth auth.Auth, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "auth token not provided", http.StatusUnauthorized)
			return
		}

		username, err := auth.ValidateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxUsernameKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
