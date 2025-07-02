package repository

import (
	"database/sql"
	"fmt"

	"real-time-voting/internal/models"

	"github.com/google/uuid"
)

type VoteRepository struct {
	db *sql.DB
}

func NewVoteRepository(db *sql.DB) *VoteRepository {
	return &VoteRepository{db: db}
}

func (r *VoteRepository) Create(tx *sql.Tx, vote *models.Vote) error {
	query := `
		INSERT INTO votes (id, poll_id, option_id, voter_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	var err error
	if tx != nil {
		_, err = tx.Exec(query, vote.ID, vote.PollID, vote.OptionID, vote.VoterID, vote.CreatedAt)
	} else {
		_, err = r.db.Exec(query, vote.ID, vote.PollID, vote.OptionID, vote.VoterID, vote.CreatedAt)
	}
	if err != nil {
		return fmt.Errorf("failed to create vote: %w", err)
	}

	return nil
}

func (r *VoteRepository) GetVoteCountsByPoll(tx *sql.Tx, pollID uuid.UUID) (map[uuid.UUID]int, error) {
	query := `
		SELECT option_id, COUNT(*) as vote_count
		FROM votes
		WHERE poll_id = $1
		GROUP BY option_id
	`

	var rows *sql.Rows
	var err error
	if tx != nil {
		rows, err = tx.Query(query, pollID)
	} else {
		rows, err = r.db.Query(query, pollID)
	}
	if err != nil {
		return make(map[uuid.UUID]int), fmt.Errorf("failed to get vote counts: %w", err)
	}
	defer rows.Close()

	voteCounts := make(map[uuid.UUID]int)
	for rows.Next() {
		var optionID uuid.UUID
		var count int
		err := rows.Scan(&optionID, &count)
		if err != nil {
			return make(map[uuid.UUID]int), fmt.Errorf("failed to scan vote count: %w", err)
		}
		voteCounts[optionID] = count
	}

	return voteCounts, nil
}

func (r *VoteRepository) GetTotalVotesByPoll(tx *sql.Tx, pollID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM votes
		WHERE poll_id = $1
	`

	var count int
	var err error
	if tx != nil {
		err = tx.QueryRow(query, pollID).Scan(&count)
	} else {
		err = r.db.QueryRow(query, pollID).Scan(&count)
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get total votes: %w", err)
	}

	return count, nil
}

func (r *VoteRepository) HasVoted(tx *sql.Tx, pollID uuid.UUID, voterID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM votes 
			WHERE poll_id = $1 AND voter_id = $2
		)
	`

	var exists bool
	var err error
	if tx != nil {
		err = tx.QueryRow(query, pollID, voterID).Scan(&exists)
	} else {
		err = r.db.QueryRow(query, pollID, voterID).Scan(&exists)
	}
	if err != nil {
		return false, fmt.Errorf("failed to check if user voted: %w", err)
	}

	return exists, nil
}

func (r *VoteRepository) GetVotesByPoll(tx *sql.Tx, pollID uuid.UUID) ([]*models.Vote, error) {
	query := `
		SELECT id, poll_id, option_id, voter_id, created_at
		FROM votes
		WHERE poll_id = $1
		ORDER BY created_at DESC
	`

	var rows *sql.Rows
	var err error
	if tx != nil {
		rows, err = tx.Query(query, pollID)
	} else {
		rows, err = r.db.Query(query, pollID)
	}
	if err != nil {
		return []*models.Vote{}, fmt.Errorf("failed to get votes: %w", err)
	}
	defer rows.Close()

	var votes []*models.Vote
	for rows.Next() {
		vote := &models.Vote{}
		err := rows.Scan(&vote.ID, &vote.PollID, &vote.OptionID, &vote.VoterID, &vote.CreatedAt)
		if err != nil {
			return []*models.Vote{}, fmt.Errorf("failed to scan vote: %w", err)
		}
		votes = append(votes, vote)
	}

	return votes, nil
}
