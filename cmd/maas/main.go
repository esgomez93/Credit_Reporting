package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"maas/internal/config"
	"maas/internal/store"
	"maas/pkg/api"
	"maas/pkg/repository"
	"maas/pkg/service"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.yaml") // Assuming you have a config.yaml file
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	// Initialize the database
	db, err := store.NewDB(cfg.Database)
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()

	// Initialize repository, service, and API handler
	memeRepo := repository.NewMemeRepository(db)
	memeService := service.NewMemeService(memeRepo)
	memeHandler := api.NewMemeHandler(memeService)

	// Set up the router and middleware
	r := mux.NewRouter()
	api.RegisterRoutes(r, memeHandler)

	// Start the server
	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
	}

	log.Printf("Starting server on port %d", cfg.Server.Port)
	log.Fatal(srv.ListenAndServe())
}
