package dto

import (
	"time"

	"github.com/google/uuid"
)

type PollResponse struct {
	ID         uuid.UUID      `json:"id"`
	Question   string         `json:"question"`
	Options    []OptionResult `json:"options"`
	IsActive   bool           `json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	EndsAt     *time.Time     `json:"ends_at,omitempty"`
	CreatedBy  *uuid.UUID     `json:"created_by,omitempty"`
	TotalVotes int            `json:"total_votes"`
}

type OptionResult struct {
	ID         uuid.UUID `json:"id"`
	Text       string    `json:"text"`
	VoteCount  int       `json:"vote_count"`
	Percentage float64   `json:"percentage"`
}

type PollUpdateEvent struct {
	Type string       `json:"type"`
	Poll PollResponse `json:"poll"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	User      User   `json:"user"`
	ExpiresAt int64  `json:"expires_at"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
