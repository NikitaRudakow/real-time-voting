package repository

import (
	"database/sql"
	"fmt"
	"time"

	"real-time-voting/internal/models"

	"github.com/google/uuid"
)

type PollRepository struct {
	db *sql.DB
}

func NewPollRepository(db *sql.DB) *PollRepository {
	return &PollRepository{db: db}
}

func (r *PollRepository) Create(tx *sql.Tx, poll *models.Poll, options []string) error {
	// Insert poll
	pollQuery := `
		INSERT INTO polls (id, question, is_active, created_at, updated_at, ends_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := tx.Exec(pollQuery, poll.ID, poll.Question, poll.IsActive, poll.CreatedAt, poll.UpdatedAt, poll.EndsAt, poll.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to insert poll: %w", err)
	}

	// Insert options
	optionQuery := `
		INSERT INTO poll_options (id, poll_id, text)
		VALUES ($1, $2, $3)
	`
	for _, optionText := range options {
		optionID := uuid.New()
		_, err = tx.Exec(optionQuery, optionID, poll.ID, optionText)
		if err != nil {
			return fmt.Errorf("failed to insert option: %w", err)
		}
	}

	return nil
}

func (r *PollRepository) GetByID(tx *sql.Tx, id uuid.UUID) (*models.Poll, error) {
	query := `
		SELECT p.id, p.question, p.is_active, p.created_at, p.updated_at, p.ends_at, p.created_by,
		       o.id, o.poll_id, o.text,
		       COALESCE(COUNT(v.id), 0) as vote_count
		FROM polls p
		LEFT JOIN poll_options o ON p.id = o.poll_id
		LEFT JOIN votes v ON o.id = v.option_id
		WHERE p.id = $1
		GROUP BY p.id, p.question, p.is_active, p.created_at, p.updated_at, p.ends_at, p.created_by, o.id, o.poll_id, o.text
		ORDER BY o.id
	`

	var rows *sql.Rows
	var err error
	if tx != nil {
		rows, err = tx.Query(query, id)
	} else {
		rows, err = r.db.Query(query, id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get poll with options: %w", err)
	}
	defer rows.Close()

	var poll *models.Poll
	var options []models.Option
	var totalVotes int
	for rows.Next() {
		var (
			pollID               uuid.UUID
			question             string
			isActive             bool
			createdAt, updatedAt time.Time
			endsAt               sql.NullTime
			createdBy            uuid.UUID
			optionID             sql.NullString
			optionPollID         sql.NullString
			optionText           sql.NullString
			voteCount            int
		)

		err := rows.Scan(
			&pollID, &question, &isActive, &createdAt, &updatedAt, &endsAt, &createdBy,
			&optionID, &optionPollID, &optionText,
			&voteCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan poll with options: %w", err)
		}

		if poll == nil {
			createdByPtr := createdBy
			poll = &models.Poll{
				ID:        pollID,
				Question:  question,
				IsActive:  isActive,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				EndsAt:    nil,
				CreatedBy: &createdByPtr,
			}
			if endsAt.Valid {
				poll.EndsAt = &endsAt.Time
			}
		}

		if optionID.Valid && optionPollID.Valid && optionText.Valid {
			optID, err1 := uuid.Parse(optionID.String)
			optPollID, err2 := uuid.Parse(optionPollID.String)
			if err1 == nil && err2 == nil {
				options = append(options, models.Option{
					ID:        optID,
					PollID:    optPollID,
					Text:      optionText.String,
					VoteCount: voteCount,
				})
				totalVotes += voteCount
			}
		}
	}

	if poll == nil {
		return nil, sql.ErrNoRows
	}
	poll.Options = options
	poll.TotalVotes = totalVotes
	return poll, nil
}

func (r *PollRepository) GetAll(tx *sql.Tx) ([]*models.Poll, error) {
	query := `
		SELECT p.id, p.question, p.is_active, p.created_at, p.updated_at, p.ends_at, p.created_by,
		       o.id, o.poll_id, o.text,
		       COALESCE(COUNT(v.id), 0) as vote_count
		FROM polls p
		LEFT JOIN poll_options o ON p.id = o.poll_id
		LEFT JOIN votes v ON o.id = v.option_id
		GROUP BY p.id, p.question, p.is_active, p.created_at, p.updated_at, p.ends_at, p.created_by, o.id, o.poll_id, o.text
		ORDER BY p.created_at DESC, o.id
	`

	var rows *sql.Rows
	var err error
	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = r.db.Query(query)
	}
	if err != nil {
		return []*models.Poll{}, fmt.Errorf("failed to get polls: %w", err)
	}
	defer rows.Close()

	pollMap := make(map[uuid.UUID]*models.Poll)

	for rows.Next() {
		var (
			pollID               uuid.UUID
			question             string
			isActive             bool
			createdAt, updatedAt time.Time
			endsAt               sql.NullTime
			createdBy            uuid.UUID
			optionID             sql.NullString
			optionPollID         sql.NullString
			optionText           sql.NullString
			voteCount            int
		)

		err := rows.Scan(
			&pollID, &question, &isActive, &createdAt, &updatedAt, &endsAt, &createdBy,
			&optionID, &optionPollID, &optionText,
			&voteCount,
		)
		if err != nil {
			return []*models.Poll{}, fmt.Errorf("failed to scan poll with options: %w", err)
		}

		poll, exists := pollMap[pollID]
		if !exists {
			createdByPtr := createdBy
			poll = &models.Poll{
				ID:        pollID,
				Question:  question,
				IsActive:  isActive,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				EndsAt:    nil,
				CreatedBy: &createdByPtr,
			}
			if endsAt.Valid {
				poll.EndsAt = &endsAt.Time
			}
			poll.Options = []models.Option{}
			pollMap[pollID] = poll
		}

		if optionID.Valid && optionPollID.Valid && optionText.Valid {
			optID, err1 := uuid.Parse(optionID.String)
			optPollID, err2 := uuid.Parse(optionPollID.String)
			if err1 == nil && err2 == nil {
				poll.Options = append(poll.Options, models.Option{
					ID:        optID,
					PollID:    optPollID,
					Text:      optionText.String,
					VoteCount: voteCount,
				})
				poll.TotalVotes += voteCount
			}
		}
	}

	var polls []*models.Poll
	for _, poll := range pollMap {
		polls = append(polls, poll)
	}

	return polls, nil
}

func (r *PollRepository) GetActive(tx *sql.Tx) ([]*models.Poll, error) {
	query := `
		SELECT p.id, p.question, p.is_active, p.created_at, p.updated_at, p.ends_at, p.created_by,
		       o.id, o.poll_id, o.text,
		       COALESCE(COUNT(v.id), 0) as vote_count
		FROM polls p
		LEFT JOIN poll_options o ON p.id = o.poll_id
		LEFT JOIN votes v ON o.id = v.option_id
		WHERE p.is_active = true AND (p.ends_at IS NULL OR p.ends_at > $1)
		GROUP BY p.id, p.question, p.is_active, p.created_at, p.updated_at, p.ends_at, p.created_by, o.id, o.poll_id, o.text
		ORDER BY p.created_at DESC, o.id
	`

	var rows *sql.Rows
	var err error
	if tx != nil {
		rows, err = tx.Query(query, time.Now())
	} else {
		rows, err = r.db.Query(query, time.Now())
	}
	if err != nil {
		return []*models.Poll{}, fmt.Errorf("failed to get active polls: %w", err)
	}
	defer rows.Close()

	pollMap := make(map[uuid.UUID]*models.Poll)

	for rows.Next() {
		var (
			pollID               uuid.UUID
			question             string
			isActive             bool
			createdAt, updatedAt time.Time
			endsAt               sql.NullTime
			createdBy            uuid.UUID
			optionID             sql.NullString
			optionPollID         sql.NullString
			optionText           sql.NullString
			voteCount            int
		)

		err := rows.Scan(
			&pollID, &question, &isActive, &createdAt, &updatedAt, &endsAt, &createdBy,
			&optionID, &optionPollID, &optionText,
			&voteCount,
		)
		if err != nil {
			return []*models.Poll{}, fmt.Errorf("failed to scan poll with options: %w", err)
		}

		poll, exists := pollMap[pollID]
		if !exists {
			createdByPtr := createdBy
			poll = &models.Poll{
				ID:        pollID,
				Question:  question,
				IsActive:  isActive,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				EndsAt:    nil,
				CreatedBy: &createdByPtr,
			}
			if endsAt.Valid {
				poll.EndsAt = &endsAt.Time
			}
			poll.Options = []models.Option{}
			pollMap[pollID] = poll
		}

		if optionID.Valid && optionPollID.Valid && optionText.Valid {
			optID, err1 := uuid.Parse(optionID.String)
			optPollID, err2 := uuid.Parse(optionPollID.String)
			if err1 == nil && err2 == nil {
				poll.Options = append(poll.Options, models.Option{
					ID:        optID,
					PollID:    optPollID,
					Text:      optionText.String,
					VoteCount: voteCount,
				})
				poll.TotalVotes += voteCount
			}
		}
	}

	var polls []*models.Poll
	for _, poll := range pollMap {
		polls = append(polls, poll)
	}

	return polls, nil
}

func (r *PollRepository) Update(tx *sql.Tx, poll *models.Poll) error {
	query := `
		UPDATE polls
		SET question = $2, is_active = $3, updated_at = $4, ends_at = $5
		WHERE id = $1
	`

	_, err := tx.Exec(query, poll.ID, poll.Question, poll.IsActive, poll.UpdatedAt, poll.EndsAt)
	if err != nil {
		return fmt.Errorf("failed to update poll: %w", err)
	}

	return nil
}

func (r *PollRepository) Delete(tx *sql.Tx, id uuid.UUID) error {
	_, err := tx.Exec(`DELETE FROM votes WHERE poll_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete votes: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM poll_options WHERE poll_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete poll options: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM polls WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete poll: %w", err)
	}

	return nil
}
