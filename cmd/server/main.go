package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"real-time-voting/internal/api"
	"real-time-voting/internal/config"
	"real-time-voting/internal/database"
	"real-time-voting/internal/logger"
	"real-time-voting/internal/repository"
	"real-time-voting/internal/service"
	"real-time-voting/internal/websocket"

	_ "real-time-voting/docs" // Swagger docs
)

// @title Real-Time Voting API
// @version 1.0
// @description Real-time polling system for "The Decisives" debate club
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token only

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := logger.New(cfg.LogLevel)

	// Initialize database
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		logger.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	// Initialize repositories
	pollRepo := repository.NewPollRepository(db)
	voteRepo := repository.NewVoteRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Initialize WebSocket hub
	wsHub := websocket.NewHub(logger)
	go wsHub.Run()

	// Initialize services
	pollService := service.NewPollService(pollRepo, voteRepo, wsHub, db, logger)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, db, logger)

	// Initialize API server
	server := api.NewServer(cfg, pollService, authService, wsHub, logger)

	// Start server in a goroutine
	go func() {
		logger.Info().Str("port", cfg.Port).Msg("Starting server")
		if err := server.Start(); err != nil {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server exited")
}
