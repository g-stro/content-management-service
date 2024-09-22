package main

import (
	"github.com/g-stro/content-service/internal/database"
	"github.com/g-stro/content-service/internal/domain/content"
	"github.com/g-stro/content-service/internal/domain/content/repository"
	"github.com/g-stro/content-service/middleware"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// Load configs
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8080"
	}
	// Connect to database
	conn := database.NewConnection()
	defer conn.Close()
	// Create multiplexer/router
	mux := http.NewServeMux()
	// Create repositories
	contentRepo := repository.NewPostgresContentRepository(conn)
	// Register services
	content.NewContentService(mux, contentRepo)
	// Setup middleware
	handler := middleware.CorsMiddleware(mux)
	// Create HTTP server
	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
