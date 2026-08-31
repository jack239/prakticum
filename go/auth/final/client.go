package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secretKey = []byte("qw45jk32")
)

func GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func main() {
	token, err := GenerateToken(123)
	if err != nil {
		fmt.Println("Generation token:", err)
		return
	}

	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/home", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	request.Header.Set("Authorization", token)
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(bodyBytes))
	}
}
