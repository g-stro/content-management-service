package main

import (
	"github.com/g-stro/content-service/database"
	"github.com/g-stro/content-service/internal/http/handler"
	"github.com/g-stro/content-service/internal/http/middleware"
	"github.com/g-stro/content-service/internal/repository"
	"github.com/g-stro/content-service/internal/service"
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
	conn, err := database.NewConnection()
	if err != nil {
		slog.Error("failed to establish database connection", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Create repository
	contentRepo := repository.NewPostgresContentRepository(conn)
	// Create service
	contentService := service.NewContentService(contentRepo)
	// Create handler
	contentHandler := handler.NewContentHandler(contentService)

	// Create multiplexer (router)
	mux := http.NewServeMux()
	// Register routes
	contentHandler.RegisterRoutes(mux)
	// Setup middleware
	httpHandler := middleware.CorsMiddleware(mux)

	// Create HTTP server
	err = http.ListenAndServe(":"+port, httpHandler)
	if err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
