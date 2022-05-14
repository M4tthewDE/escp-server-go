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
	http.HandleFunc("/ranking", handler.HandleRanking)
	http.HandleFunc("/lock", handler.HandleLock)
	http.HandleFunc("/done", handler.HandleDone)
	http.HandleFunc("/result", handler.GetResult)

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
