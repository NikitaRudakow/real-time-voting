package api

import (
	"net/http"

	"real-time-voting/internal/dto"

	"github.com/gin-gonic/gin"
)

// @Summary Register new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "User registration data"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /auth/register [post]
func (s *Server) register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.handleValidationError(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		s.handleError(c, err)
		return
	}

	user, err := s.authService.Register(&req)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to register user")
		s.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

// @Summary Login user
// @Description Authenticate user and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (s *Server) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.handleValidationError(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		s.handleError(c, err)
		return
	}

	response, err := s.authService.Login(&req)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to login user")
		s.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
