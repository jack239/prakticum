package main

import (
	"log"
	"net/http"

	"utils/server_utils"
)

func main() {
	http.HandleFunc("/home", server_utils.AuthMiddleware(server_utils.HandlePage))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Fail", err)
	}

}
