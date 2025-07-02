package repository

import (
	"database/sql"
	"fmt"

	"real-time-voting/internal/models"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(tx *sql.Tx, user *models.User) error {
	query := `
		INSERT INTO users (id, username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	var err error
	if tx != nil {
		_, err = tx.Exec(query, user.ID, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	} else {
		_, err = r.db.Exec(query, user.ID, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	}
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(tx *sql.Tx, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	var err error
	if tx != nil {
		err = tx.QueryRow(query, id).Scan(
			&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
		)
	} else {
		err = r.db.QueryRow(query, id).Scan(
			&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByUsername(tx *sql.Tx, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	user := &models.User{}
	var err error
	if tx != nil {
		err = tx.QueryRow(query, username).Scan(
			&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
		)
	} else {
		err = r.db.QueryRow(query, username).Scan(
			&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(tx *sql.Tx, email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	var err error
	if tx != nil {
		err = tx.QueryRow(query, email).Scan(
			&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
		)
	} else {
		err = r.db.QueryRow(query, email).Scan(
			&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *UserRepository) Update(tx *sql.Tx, user *models.User) error {
	query := `
		UPDATE users
		SET username = $2, email = $3, password = $4, updated_at = $5
		WHERE id = $1
	`

	var err error
	if tx != nil {
		_, err = tx.Exec(query, user.ID, user.Username, user.Email, user.Password, user.UpdatedAt)
	} else {
		_, err = r.db.Exec(query, user.ID, user.Username, user.Email, user.Password, user.UpdatedAt)
	}
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *UserRepository) Delete(tx *sql.Tx, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	var err error
	if tx != nil {
		_, err = tx.Exec(query, id)
	} else {
		_, err = r.db.Exec(query, id)
	}
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetAll(tx *sql.Tx) ([]*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	var rows *sql.Rows
	var err error
	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = r.db.Query(query)
	}
	if err != nil {
		return []*models.User{}, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return []*models.User{}, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}
