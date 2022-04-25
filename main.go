package main

import (
	"log"
	"net/http"
	"os"

	"github.com/m4tthewde/escp-server-go/internal/api"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Print("starting server...")

	handler := api.NewHandler()

	http.HandleFunc("/countries", handler.GetCountries)
	http.HandleFunc("/result", handler.SetResult)
	http.HandleFunc("/ranking", handler.SetRanking)
	http.HandleFunc("/user", handler.HandleUser)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
