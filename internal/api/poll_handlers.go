package api

import (
	"net/http"

	"real-time-voting/internal/dto"
	customErrors "real-time-voting/internal/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Create a new poll
// @Description Create a new poll with multiple choice options
// @Tags polls
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param poll body dto.CreatePollRequest true "Poll creation data"
// @Success 201 {object} dto.PollResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /polls [post]
func (s *Server) createPoll(c *gin.Context) {
	var req dto.CreatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.handleValidationError(c, err)
		return
	}

	// Валидация входных данных
	if err := req.Validate(); err != nil {
		s.handleError(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		s.handleError(c, customErrors.ErrUnauthorized)
		return
	}

	poll, err := s.pollService.CreatePoll(&req, userID.(uuid.UUID))
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create poll")
		s.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, poll)
}

// @Summary Get all polls
// @Description Get list of all polls
// @Tags polls
// @Accept json
// @Produce json
// @Param page query int false "Page number" minimum(1)
// @Param page_size query int false "Page size" minimum(1) maximum(100)
// @Success 200 {array} dto.PollResponse
// @Router /polls [get]
func (s *Server) getAllPolls(c *gin.Context) {
	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		s.handleValidationError(c, err)
		return
	}
	if pagination.Page == 0 {
		pagination.Page = 1
	}
	if pagination.PageSize == 0 {
		pagination.PageSize = 10
	}

	polls, err := s.pollService.GetPolls(pagination.Page, pagination.PageSize, false)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get all polls")
		s.handleDatabaseError(c, err)
		return
	}

	if polls == nil {
		polls = []dto.PollResponse{}
	}

	c.JSON(http.StatusOK, polls)
}

// @Summary Get active polls
// @Description Get list of currently active polls
// @Tags polls
// @Accept json
// @Produce json
// @Param page query int false "Page number" minimum(1)
// @Param page_size query int false "Page size" minimum(1) maximum(100)
// @Success 200 {array} dto.PollResponse
// @Router /polls/active [get]
func (s *Server) getActivePolls(c *gin.Context) {
	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		s.handleValidationError(c, err)
		return
	}
	if pagination.Page == 0 {
		pagination.Page = 1
	}
	if pagination.PageSize == 0 {
		pagination.PageSize = 10
	}

	polls, err := s.pollService.GetPolls(pagination.Page, pagination.PageSize, true)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get active polls")
		s.handleDatabaseError(c, err)
		return
	}

	if polls == nil {
		polls = []dto.PollResponse{}
	}

	c.JSON(http.StatusOK, polls)
}

// @Summary Get poll by ID
// @Description Get specific poll with results
// @Tags polls
// @Accept json
// @Produce json
// @Param id path string true "Poll ID"
// @Success 200 {object} dto.PollResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /polls/{id} [get]
func (s *Server) getPoll(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		s.handleValidationError(c, err)
		return
	}

	poll, err := s.pollService.GetPoll(id)
	if err != nil {
		s.logger.Error().Err(err).Str("poll_id", id.String()).Msg("Failed to get poll")
		s.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, poll)
}

// @Summary Update poll
// @Description Update an existing poll
// @Tags polls
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Poll ID"
// @Param poll body dto.CreatePollRequest true "Updated poll data"
// @Success 200 {object} dto.PollResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /polls/{id} [put]
func (s *Server) updatePoll(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		s.handleValidationError(c, err)
		return
	}

	var req dto.CreatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.handleValidationError(c, err)
		return
	}

	// Валидация входных данных
	if err := req.Validate(); err != nil {
		s.handleError(c, err)
		return
	}

	poll, err := s.pollService.UpdatePoll(id, &req)
	if err != nil {
		s.logger.Error().Err(err).Str("poll_id", id.String()).Msg("Failed to update poll")
		s.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, poll)
}

// @Summary Delete poll
// @Description Delete an existing poll
// @Tags polls
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Poll ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /polls/{id} [delete]
func (s *Server) deletePoll(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		s.handleValidationError(c, err)
		return
	}

	err = s.pollService.DeletePoll(id)
	if err != nil {
		s.logger.Error().Err(err).Str("poll_id", id.String()).Msg("Failed to delete poll")
		s.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Poll deleted successfully"})
}

// @Summary Cast a vote
// @Description Vote on a specific poll option
// @Tags polls
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Poll ID"
// @Param vote body dto.VoteRequest true "Vote data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /polls/{id}/vote [post]
func (s *Server) vote(c *gin.Context) {
	pollID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		s.handleValidationError(c, err)
		return
	}

	var req dto.VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.handleValidationError(c, err)
		return
	}

	// Валидация входных данных
	if err := req.Validate(); err != nil {
		s.handleError(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		s.handleError(c, customErrors.ErrUnauthorized)
		return
	}

	err = s.pollService.Vote(pollID, req.OptionID, userID.(uuid.UUID))
	if err != nil {
		s.logger.Error().Err(err).Str("poll_id", pollID.String()).Msg("Failed to cast vote")
		s.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote cast successfully"})
}
