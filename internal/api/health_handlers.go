package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Health check
// @Description Get server health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":            "ok",
		"message":           "Real-time Voting API is running",
		"websocket_clients": s.wsHub.GetClientCount(),
	})
}
