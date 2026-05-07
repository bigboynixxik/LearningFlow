package service

import (
	"context"
	"learningflow/internal/models"
)

type SubjectRepository interface {
	GetAll(ctx context.Context) ([]models.Subject, error)
	GetByID(ctx context.Context, id int64) (*models.Subject, error)
}

type SubjectService struct {
	repo SubjectRepository
}

func NewSubjectService(repo SubjectRepository) *SubjectService {
	return &SubjectService{repo: repo}
}

func (s *SubjectService) GetAll(ctx context.Context) ([]models.Subject, error) {
	return s.repo.GetAll(ctx)
}

func (s *SubjectService) GetByID(ctx context.Context, id int64) (*models.Subject, error) {
	return s.repo.GetByID(ctx, id)
}
