package service

import (
	"context"
	"errors"
	"fmt"
	"learningflow/internal/models"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	GetByToken(ctx context.Context, token string) (*models.Session, error)
	Delete(ctx context.Context, token string) error // Для Logout
}

type AuthService struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
}

func NewAuthService(ur UserRepository, sr SessionRepository) *AuthService {
	return &AuthService{
		userRepo:    ur,
		sessionRepo: sr}
}

func (s *AuthService) Register(ctx context.Context, email, plainPassword string, role models.Role) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("service.register hash error: %w", err)
	}

	var user models.User
	user.Email = email
	user.Role = role
	user.PasswordHash = string(hash)

	if err := s.userRepo.Create(ctx, &user); err != nil {
		if errors.Is(err, models.ErrAlreadyExists) {
			return "", models.ErrAlreadyExists
		}
		return "", fmt.Errorf("service.register create user error: %w", err)
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(24 * time.Hour)
	session := &models.Session{
		Token:     sessionToken,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return "", fmt.Errorf("service.register create session error: %w", err)
	}

	return sessionToken, nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return "", fmt.Errorf("invalid credentials")
		}
		return "", fmt.Errorf("service.login: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("service.login: %w", err)
	}
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(24 * time.Hour)
	session := &models.Session{
		Token:     sessionToken,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}
	err = s.sessionRepo.Create(ctx, session)
	if err != nil {
		return "", fmt.Errorf("service.login: %w", err)
	}
	return sessionToken, nil
}

func (s *AuthService) ValidateSession(ctx context.Context, token string) (string, error) {
	session, err := s.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		return "", fmt.Errorf("session not found: %w", err)
	}

	if time.Now().After(session.ExpiresAt) {
		_ = s.sessionRepo.Delete(ctx, token)
		return "", fmt.Errorf("session expired")
	}

	return session.UserID, nil
}
