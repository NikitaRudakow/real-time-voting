package api

import (
	"context"
	"net/http"

	"real-time-voting/internal/config"
	"real-time-voting/internal/service"
	"real-time-voting/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config      *config.Config
	pollService *service.PollService
	authService *service.AuthService
	wsHub       *websocket.Hub
	logger      zerolog.Logger
	server      *http.Server
}

func NewServer(cfg *config.Config, pollService *service.PollService, authService *service.AuthService, wsHub *websocket.Hub, logger zerolog.Logger) *Server {
	return &Server{
		config:      cfg,
		pollService: pollService,
		authService: authService,
		wsHub:       wsHub,
		logger:      logger,
	}
}

func (s *Server) Start() error {
	router := gin.Default()

	// Add middleware
	router.Use(s.LoggingMiddleware())
	router.Use(s.ErrorHandlingMiddleware())

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", s.register)
			auth.POST("/login", s.login)
		}

		// Health check
		api.GET("/health", s.healthCheck)

		// Poll routes
		polls := api.Group("/polls")
		{
			polls.GET("/", s.getAllPolls)          // Public
			polls.GET("/active", s.getActivePolls) // Public
			polls.GET("/:id", s.getPoll)           // Public

			// Protected routes
			protected := polls.Group("/")
			protected.Use(s.AuthMiddleware())
			{
				protected.POST("/", s.createPoll)
				protected.PUT("/:id", s.updatePoll)
				protected.DELETE("/:id", s.deletePoll)
				protected.POST("/:id/vote", s.vote)
			}
		}

		// WebSocket route
		api.GET("/ws", s.handleWebSocket)
	}

	// Create HTTP server
	s.server = &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: router,
	}

	s.logger.Info().Str("port", s.config.Port).Msg("Starting HTTP server")
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) handleWebSocket(c *gin.Context) {
	websocket.ServeWs(s.wsHub, c.Writer, c.Request)
}
