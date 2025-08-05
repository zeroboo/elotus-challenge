package main

import (
	"net/http"
	"os"

	"elotuschallenge/database"
	"elotuschallenge/handler"
	"elotuschallenge/middleware"

	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer database.CloseDB()

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set up routes
	setupRoutes()

	log.Info().Str("port", port).Msg("Server is starting...")

	// Start server
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

// setupRoutes initializes the HTTP routes for the server using net/http
func setupRoutes() {
	// Authentication routes
	http.HandleFunc("/register", handler.HandleRegister)
	http.HandleFunc("/login", handler.HandleLogin)

	// Protected section
	http.HandleFunc("/upload", middleware.AuthUser(handler.HandleUpload))

	// Health check
	http.HandleFunc("/health", handler.HandleHealth)

	log.Info().Msg("Routes configured")
}
