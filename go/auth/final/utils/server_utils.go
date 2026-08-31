package utils

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secret = []byte("qw45jk32")
)

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}

type Message struct {
	Status string `json:"status"`
	Info   string `json:"info"`
}

func HandlePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("you have gained access"))
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "token not found", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := ParseToken(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "wrong token", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
