package main

import (
	"net/http"
	"os"

	"elotuschallenge/database"
	"elotuschallenge/handler"
	"elotuschallenge/internal"
	"elotuschallenge/middleware"

	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer database.CloseDB()

	// Initialize services
	internal.InitServices()

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
	// API routes
	http.HandleFunc("/api/register", handler.HandleRegister)
	http.HandleFunc("/api/login", handler.HandleLogin)
	http.HandleFunc("/api/upload", middleware.AuthUser(handler.HandleUpload))

	// Form routes (static files)
	http.HandleFunc("/form/register", handler.HandleStatic)
	http.HandleFunc("/form/login", handler.HandleStatic)
	http.HandleFunc("/form/upload", handler.HandleStatic)

	// Health check
	http.HandleFunc("/health", handler.HandleHealth)

	log.Info().Msg("Routes configured")
}
