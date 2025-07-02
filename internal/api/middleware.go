package api

import (
	"bytes"
	"io"
	"strings"
	"time"

	customErrors "real-time-voting/internal/errors"

	"github.com/gin-gonic/gin"
)

func (s *Server) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			s.handleError(c, customErrors.ErrUnauthorized)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			s.handleError(c, customErrors.ErrUnauthorized)
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := s.authService.ValidateToken(token)
		if err != nil {
			s.handleError(c, customErrors.ErrUnauthorized)
			c.Abort()
			return
		}

		// Добавляем claims в контекст
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("claims", claims)

		c.Next()
	}
}

func (s *Server) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]
		claims, err := s.authService.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		// Добавляем claims в контекст
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("claims", claims)

		c.Next()
	}
}

func (s *Server) LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		userID, _ := c.Get("user_id")
		query := c.Request.URL.RawQuery

		var requestBody string
		if method == "POST" || method == "PUT" || method == "PATCH" || method == "DELETE" {
			if c.Request.Body != nil {
				bodyBytes, err := io.ReadAll(c.Request.Body)
				if err == nil {
					if len(bodyBytes) > 0 && len(bodyBytes) < 4096 {
						requestBody = string(bodyBytes)
					}
					c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				}
			}
		}

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logEvent := s.logger.Info().
			Str("method", method).
			Str("path", path).
			Int("status", status).
			Dur("latency", latency).
			Interface("user_id", userID).
			Str("query", query)

		if requestBody != "" {
			logEvent = logEvent.Str("body", requestBody)
		}

		logEvent.Msg("HTTP request")
	}
}

func (s *Server) ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				s.logger.Error().Interface("panic", rec).Str("path", c.Request.URL.Path).Msg("Panic recovered")
				c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error"})
			}
		}()
		c.Next()
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				s.logger.Error().Err(e.Err).Str("path", c.Request.URL.Path).Msg("API error")
			}
		}
	}
}
