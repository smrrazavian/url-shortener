package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/smrrazavian/url-shortener/internal/config"
	"github.com/smrrazavian/url-shortener/internal/router"
	"github.com/smrrazavian/url-shortener/internal/router/handlers"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("Starting server...")
	if err := handlers.LoadFromFile("data.json"); err != nil {
		log.Fatal("Error loading store:", err)
	}

	mux := http.NewServeMux()
	router.AddRoutes(mux)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: mux,
	}

	log.Println("Listened on port", cfg.ServerPort)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT)
	<-stop

	log.Println("Shutting Down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Saving URL store to file...")
	if err := handlers.StoreToFile("data.json"); err != nil {
		log.Panicln("Error during save", err)
	}
	log.Println("Shutdown complete.")
}
