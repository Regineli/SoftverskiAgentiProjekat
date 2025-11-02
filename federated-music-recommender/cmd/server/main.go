package main

import (
	"fmt"
	"log"
	"net/http"

	"federated-music-recommender/internal/server"
)

func main() {
	fmt.Println("🎵 Federated Server started on port 8080...")

	http.HandleFunc("/upload", server.HandleUpload)
	http.HandleFunc("/aggregate", server.HandleAggregate)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
