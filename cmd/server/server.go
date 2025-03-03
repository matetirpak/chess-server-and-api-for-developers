package main

import (
	"log"
	"net/http"

	"github.com/matetirpak/chess-server-and-api-for-developers/pkg/server"
)

func main() {
	log.Printf("Server started")

	router := server.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
