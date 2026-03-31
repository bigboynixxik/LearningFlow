package service

import (
	"context"
	"fmt"
	"learningflow/internal/domain"
	"learningflow/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    domain.UserRepository
	sessionRepo domain.SessionRepository
}

func NewAuthService(ur domain.UserRepository, sr domain.SessionRepository) *AuthService {
	return &AuthService{
		userRepo:    ur,
		sessionRepo: sr}
}

func (s *AuthService) Register(ctx context.Context, email, plainPassword string, role models.Role) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("service.register: %w", err)
	}
	var user models.User
	user.Email = email
	user.Role = role
	user.PasswordHash = string(hash)
	if err := s.userRepo.Create(ctx, &user); err != nil {
		return nil, fmt.Errorf("service.register: %w", err)
	}
	return &user, nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("service.login: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("service.login: %w", err)
	}
	sessionToken := uuid.NewString()
	return sessionToken, nil
}
