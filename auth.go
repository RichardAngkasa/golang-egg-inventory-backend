package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	secretKey []byte
}

func NewAuth(secret string) *Auth {
	return &Auth{secretKey: []byte(secret)}
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (a *Auth) GenerateJWT(username string) (string, error) {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secretKey)
}

func (a *Auth) ValidateJWT(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	token, err := jwt.ParseWithClaims(parts[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected sigining method")
		}
		return a.secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Username, nil
	}
	return "", errors.New("invalid token")
}
