package api

import (
	"net/http"
	customErrors "real-time-voting/internal/errors"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

func (s *Server) handleError(c *gin.Context, err error) {
	var statusCode int
	var errorCode string

	switch err {
	case customErrors.ErrInvalidEmail,
		customErrors.ErrInvalidUsername,
		customErrors.ErrInvalidPassword,
		customErrors.ErrEmptyQuestion,
		customErrors.ErrEmptyOptions,
		customErrors.ErrTooManyOptions,
		customErrors.ErrInvalidEndDate,
		customErrors.ErrEmptyOptionText:
		statusCode = http.StatusBadRequest
		errorCode = "VALIDATION_ERROR"

	case customErrors.ErrUserNotFound,
		customErrors.ErrPollNotFound,
		customErrors.ErrOptionNotFound:
		statusCode = http.StatusNotFound
		errorCode = "NOT_FOUND"

	case customErrors.ErrUnauthorized:
		statusCode = http.StatusUnauthorized
		errorCode = "UNAUTHORIZED"

	case customErrors.ErrForbidden:
		statusCode = http.StatusForbidden
		errorCode = "FORBIDDEN"

	case customErrors.ErrUserExists,
		customErrors.ErrEmailExists,
		customErrors.ErrUsernameExists:
		statusCode = http.StatusConflict
		errorCode = "CONFLICT"

	case customErrors.ErrPollNotActive,
		customErrors.ErrPollEnded,
		customErrors.ErrAlreadyVoted,
		customErrors.ErrInvalidOption:
		statusCode = http.StatusBadRequest
		errorCode = "BUSINESS_RULE_VIOLATION"

	default:
		statusCode = http.StatusInternalServerError
		errorCode = "INTERNAL_ERROR"
	}

	response := ErrorResponse{
		Error: err.Error(),
		Code:  errorCode,
	}

	c.JSON(statusCode, response)
}

func (s *Server) handleValidationError(c *gin.Context, err error) {
	response := ErrorResponse{
		Error: err.Error(),
		Code:  "VALIDATION_ERROR",
	}
	c.JSON(http.StatusBadRequest, response)
}

// handleDatabaseError обрабатывает ошибки базы данных
func (s *Server) handleDatabaseError(c *gin.Context, err error) {
	s.logger.Error().Err(err).Msg("Database error occurred")
	response := ErrorResponse{
		Error: "Internal server error",
		Code:  "DATABASE_ERROR",
	}
	c.JSON(http.StatusInternalServerError, response)
}
