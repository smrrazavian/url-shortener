package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/smrrazavian/url-shortener/internal/config"
	"github.com/smrrazavian/url-shortener/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Println("Starting server...")

	mux := http.NewServeMux()
	router.AddRoutes(mux)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: mux,
	}

	fmt.Println("Listened on port", cfg.ServerPort)

	server.ListenAndServe()
}
