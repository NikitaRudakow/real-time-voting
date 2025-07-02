package dto

import (
	customErrors "real-time-voting/internal/errors"
	"time"

	"github.com/google/uuid"
)

type CreatePollRequest struct {
	Question string     `json:"question" binding:"required"`
	Options  []string   `json:"options" binding:"required,min=2,max=10"`
	EndsAt   *time.Time `json:"ends_at,omitempty"`
}

func (r *CreatePollRequest) Validate() error {
	if r.Question == "" {
		return customErrors.ErrEmptyQuestion
	}
	if len(r.Question) > 500 {
		return customErrors.ErrEmptyQuestion
	}
	if len(r.Options) < 2 {
		return customErrors.ErrEmptyOptions
	}
	if len(r.Options) > 10 {
		return customErrors.ErrTooManyOptions
	}

	for _, option := range r.Options {
		if option == "" {
			return customErrors.ErrEmptyOptionText
		}
		if len(option) > 200 {
			return customErrors.ErrEmptyOptionText
		}
	}

	if r.EndsAt != nil && r.EndsAt.Before(time.Now()) {
		return customErrors.ErrInvalidEndDate
	}

	return nil
}

type VoteRequest struct {
	OptionID uuid.UUID `json:"option_id" binding:"required"`
}

func (r *VoteRequest) Validate() error {
	if r.OptionID == uuid.Nil {
		return customErrors.ErrInvalidOption
	}
	return nil
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (r *RegisterRequest) Validate() error {
	if r.Username == "" {
		return customErrors.ErrInvalidUsername
	}
	if len(r.Username) < 3 || len(r.Username) > 50 {
		return customErrors.ErrInvalidUsername
	}
	// Проверяем, что username содержит только буквы, цифры и подчеркивания
	for _, char := range r.Username {
		if !isAlphanumeric(char) && char != '_' {
			return customErrors.ErrInvalidUsername
		}
	}

	if r.Email == "" {
		return customErrors.ErrInvalidEmail
	}
	if len(r.Email) < 5 || len(r.Email) > 254 {
		return customErrors.ErrInvalidEmail
	}
	if !contains(r.Email, "@") || !contains(r.Email, ".") {
		return customErrors.ErrInvalidEmail
	}

	if r.Password == "" {
		return customErrors.ErrInvalidPassword
	}
	if len(r.Password) < 8 {
		return customErrors.ErrInvalidPassword
	}

	return nil
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (r *LoginRequest) Validate() error {
	if r.Username == "" {
		return customErrors.ErrInvalidUsername
	}
	if r.Password == "" {
		return customErrors.ErrInvalidPassword
	}
	return nil
}

func isAlphanumeric(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9')
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

type PaginationQuery struct {
	Page     int `form:"page" binding:"min=1"`
	PageSize int `form:"page_size" binding:"min=1,max=100"`
}
