package service

import (
	"errors"
	"real-time-voting/internal/database"
	"real-time-voting/internal/dto"
	customErrors "real-time-voting/internal/errors"
	"time"

	"real-time-voting/internal/models"
	"real-time-voting/internal/repository"

	"database/sql"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
	db        *sql.DB
	logger    zerolog.Logger
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string, db *sql.DB, logger zerolog.Logger) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		db:        db,
		logger:    logger,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
	existingUser, _ := s.userRepo.GetByUsername(nil, req.Username)
	if existingUser != nil {
		return nil, customErrors.ErrUsernameExists
	}

	existingUser, _ = s.userRepo.GetByEmail(nil, req.Email)
	if existingUser != nil {
		return nil, customErrors.ErrEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:        uuid.New(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer database.RollbackTx(tx)

	err = s.userRepo.Create(tx, user)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to commit transaction")
		return nil, err
	}

	return s.convertToUserResponse(user), nil
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(nil, req.Username)
	if err != nil {
		return nil, customErrors.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, customErrors.ErrUnauthorized
	}

	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:     token,
		User:      *s.convertToUser(user),
		ExpiresAt: expiresAt,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *AuthService) GetUserByID(userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(nil, userID)
	if err != nil {
		return nil, customErrors.ErrUserNotFound
	}

	return s.convertToUserResponse(user), nil
}

func (s *AuthService) GetAllUsers() ([]*dto.UserResponse, error) {
	users, err := s.userRepo.GetAll(nil)
	if err != nil {
		return []*dto.UserResponse{}, err
	}

	var responses []*dto.UserResponse
	for _, user := range users {
		response := s.convertToUserResponse(user)
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *AuthService) generateToken(user *models.User) (string, int64, error) {
	expiresAt := time.Now().Add(24 * time.Hour).Unix()

	claims := &models.Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiresAt, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "real-time-voting",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt, nil
}

func (s *AuthService) convertToUserResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (s *AuthService) convertToUser(user *models.User) *dto.User {
	return &dto.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
