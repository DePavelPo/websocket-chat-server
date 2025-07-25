package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	JwtKey string
}

func NewAuth(jwtKey string) *Auth {
	return &Auth{
		JwtKey: jwtKey,
	}
}

func (a *Auth) GenerateJWT(username, pass string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"username": username,
		"password": pass,
		"exp":      time.Now().Add(6 * time.Hour).Unix(),
	})
	return token.SignedString(a.JwtKey)
}
