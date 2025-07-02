package models

import (
	"time"

	"github.com/google/uuid"
)

type Poll struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	Question   string     `json:"question" db:"question"`
	Options    []Option   `json:"options" db:"-"`
	IsActive   bool       `json:"is_active" db:"is_active"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	EndsAt     *time.Time `json:"ends_at,omitempty" db:"ends_at"`
	CreatedBy  *uuid.UUID `json:"created_by,omitempty" db:"created_by"`
	TotalVotes int        `json:"total_votes" db:"-"`
}

type Option struct {
	ID        uuid.UUID `json:"id" db:"id"`
	PollID    uuid.UUID `json:"poll_id" db:"poll_id"`
	Text      string    `json:"text" db:"text"`
	VoteCount int       `json:"vote_count" db:"-"`
}

type Vote struct {
	ID        uuid.UUID `json:"id" db:"id"`
	PollID    uuid.UUID `json:"poll_id" db:"poll_id"`
	OptionID  uuid.UUID `json:"option_id" db:"option_id"`
	VoterID   uuid.UUID `json:"voter_id" db:"voter_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
