package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Auth interface {
	GenerateJWT(username string) (string, error)
	ValidateJWT(tokenString string) (string, error)
}

type AuthClient struct {
	jwtKey interface{}
}

func NewAuthClient(jwtKey string) Auth {
	return &AuthClient{
		jwtKey: jwtKey,
	}
}

func (a *AuthClient) GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(6 * time.Hour).Unix(),
	})
	return token.SignedString(a.jwtKey)
}

func (a *AuthClient) ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return a.jwtKey, nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", errors.New("username not found in token")
	}

	return username, nil
}
