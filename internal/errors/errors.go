package errors

import "errors"

// Validation errors
var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidUsername = errors.New("username must be 3-50 characters long and contain only letters, numbers, and underscores")
	ErrInvalidPassword = errors.New("password must be at least 8 characters long")
	ErrEmptyQuestion   = errors.New("poll question cannot be empty")
	ErrEmptyOptions    = errors.New("poll must have at least 2 options")
	ErrTooManyOptions  = errors.New("poll cannot have more than 10 options")
	ErrInvalidEndDate  = errors.New("end date must be in the future")
	ErrEmptyOptionText = errors.New("option text cannot be empty")
)

// Business logic errors
var (
	ErrUserNotFound   = errors.New("user not found")
	ErrPollNotFound   = errors.New("poll not found")
	ErrOptionNotFound = errors.New("option not found")
	ErrPollNotActive  = errors.New("poll is not active")
	ErrPollEnded      = errors.New("poll has ended")
	ErrAlreadyVoted   = errors.New("user has already voted on this poll")
	ErrInvalidOption  = errors.New("invalid option for this poll")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
)

// Database errors
var (
	ErrUserExists     = errors.New("user already exists")
	ErrEmailExists    = errors.New("email already exists")
	ErrUsernameExists = errors.New("username already exists")
)
