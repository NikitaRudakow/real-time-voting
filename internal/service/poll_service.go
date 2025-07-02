package service

import (
	"real-time-voting/internal/database"
	"real-time-voting/internal/dto"
	customErrors "real-time-voting/internal/errors"
	"time"

	"real-time-voting/internal/models"
	"real-time-voting/internal/repository"
	"real-time-voting/internal/websocket"

	"database/sql"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type PollService struct {
	pollRepo *repository.PollRepository
	voteRepo *repository.VoteRepository
	wsHub    *websocket.Hub
	db       *sql.DB
	logger   zerolog.Logger
}

func NewPollService(pollRepo *repository.PollRepository, voteRepo *repository.VoteRepository, wsHub *websocket.Hub, db *sql.DB, logger zerolog.Logger) *PollService {
	return &PollService{
		pollRepo: pollRepo,
		voteRepo: voteRepo,
		wsHub:    wsHub,
		db:       db,
		logger:   logger,
	}
}

func (s *PollService) CreatePoll(req *dto.CreatePollRequest, createdBy uuid.UUID) (dto.PollResponse, error) {
	poll := models.Poll{
		ID:        uuid.New(),
		Question:  req.Question,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		EndsAt:    req.EndsAt,
		CreatedBy: &createdBy,
	}

	tx, err := s.db.Begin()
	if err != nil {
		return dto.PollResponse{}, err
	}
	defer database.RollbackTx(tx)

	err = s.pollRepo.Create(tx, &poll, req.Options)
	if err != nil {
		return dto.PollResponse{}, err
	}

	err = tx.Commit()
	if err != nil {
		return dto.PollResponse{}, err
	}

	createdPoll, err := s.pollRepo.GetByID(nil, poll.ID)
	if err != nil {
		return dto.PollResponse{}, err
	}

	resp := pollToResponse(createdPoll)

	s.wsHub.Broadcast(dto.PollUpdateEvent{
		Type: "poll_created",
		Poll: resp,
	})

	return resp, nil
}

func (s *PollService) GetPoll(id uuid.UUID) (dto.PollResponse, error) {
	poll, err := s.pollRepo.GetByID(nil, id)
	if err != nil {
		return dto.PollResponse{}, err
	}
	resp := pollToResponse(poll)

	return resp, nil
}

func (s *PollService) GetPolls(page, pageSize int, onlyActive bool) ([]dto.PollResponse, error) {
	offset := (page - 1) * pageSize
	polls, err := s.pollRepo.GetPolls(nil, pageSize, offset, onlyActive)
	if err != nil {
		return []dto.PollResponse{}, err
	}

	var responses []dto.PollResponse
	for _, poll := range polls {
		responses = append(responses, pollToResponse(poll))
	}

	return responses, nil
}

func (s *PollService) Vote(pollID uuid.UUID, optionID uuid.UUID, voterID uuid.UUID) error {
	poll, err := s.pollRepo.GetByID(nil, pollID)
	if err != nil {
		return customErrors.ErrPollNotFound
	}

	if !poll.IsActive {
		return customErrors.ErrPollNotActive
	}

	if poll.EndsAt != nil && time.Now().After(*poll.EndsAt) {
		return customErrors.ErrPollEnded
	}

	hasVoted, err := s.voteRepo.HasVoted(nil, pollID, voterID)
	if err != nil {
		return err
	}

	if hasVoted {
		return customErrors.ErrAlreadyVoted
	}

	optionExists := false
	for _, option := range poll.Options {
		if option.ID == optionID {
			optionExists = true
			break
		}
	}

	if !optionExists {
		return customErrors.ErrInvalidOption
	}

	vote := &models.Vote{
		ID:        uuid.New(),
		PollID:    pollID,
		OptionID:  optionID,
		VoterID:   voterID,
		CreatedAt: time.Now(),
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer database.RollbackTx(tx)

	err = s.voteRepo.Create(tx, vote)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	updatedPoll, err := s.pollRepo.GetByID(nil, pollID)
	if err != nil {
		return err
	}

	resp := pollToResponse(updatedPoll)

	s.wsHub.Broadcast(dto.PollUpdateEvent{
		Type: "vote_cast",
		Poll: resp,
	})

	return nil
}

func (s *PollService) UpdatePoll(id uuid.UUID, req *dto.CreatePollRequest) (dto.PollResponse, error) {
	poll, err := s.pollRepo.GetByID(nil, id)
	if err != nil {
		return dto.PollResponse{}, customErrors.ErrPollNotFound
	}

	poll.Question = req.Question
	poll.UpdatedAt = time.Now()
	poll.EndsAt = req.EndsAt

	tx, err := s.db.Begin()
	if err != nil {
		return dto.PollResponse{}, err
	}
	defer database.RollbackTx(tx)

	err = s.pollRepo.Update(tx, poll)
	if err != nil {
		return dto.PollResponse{}, err
	}

	err = tx.Commit()
	if err != nil {
		return dto.PollResponse{}, err
	}

	updatedPoll, err := s.pollRepo.GetByID(nil, id)
	if err != nil {
		return dto.PollResponse{}, err
	}

	resp := pollToResponse(updatedPoll)

	s.wsHub.Broadcast(dto.PollUpdateEvent{
		Type: "poll_updated",
		Poll: resp,
	})

	return resp, nil
}

func (s *PollService) DeletePoll(id uuid.UUID) error {
	_, err := s.pollRepo.GetByID(nil, id)
	if err != nil {
		return customErrors.ErrPollNotFound
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer database.RollbackTx(tx)

	err = s.pollRepo.Delete(tx, id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	s.wsHub.Broadcast(dto.PollUpdateEvent{
		Type: "poll_deleted",
		Poll: dto.PollResponse{ID: id},
	})

	return nil
}

func pollToResponse(poll *models.Poll) dto.PollResponse {
	var optionResults []dto.OptionResult
	for _, option := range poll.Options {
		percentage := 0.0
		if poll.TotalVotes > 0 {
			percentage = float64(option.VoteCount) / float64(poll.TotalVotes) * 100
		}
		optionResults = append(optionResults, dto.OptionResult{
			ID:         option.ID,
			Text:       option.Text,
			VoteCount:  option.VoteCount,
			Percentage: percentage,
		})
	}
	return dto.PollResponse{
		ID:         poll.ID,
		Question:   poll.Question,
		Options:    optionResults,
		IsActive:   poll.IsActive,
		CreatedAt:  poll.CreatedAt,
		EndsAt:     poll.EndsAt,
		CreatedBy:  poll.CreatedBy,
		TotalVotes: poll.TotalVotes,
	}
}
