package service

import (
	"context"
	"learningflow/internal/models"
)

type TutorRepository interface {
	GetAll(ctx context.Context) ([]models.Tutor, error)
	GetByID(ctx context.Context, id string) (*models.Tutor, error)
	GetBySubjectID(ctx context.Context, subjectID int64) ([]models.Tutor, error)
}

type TutorService struct {
	repo TutorRepository
}

func NewTutorService(r TutorRepository) *TutorService {
	return &TutorService{repo: r}
}

func (s *TutorService) GetAll(ctx context.Context) ([]models.Tutor, error) {
	return s.repo.GetAll(ctx)
}

func (s *TutorService) GetByID(ctx context.Context, id string) (*models.Tutor, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TutorService) GetBySubjectID(ctx context.Context, subjectID int64) ([]models.Tutor, error) {
	return s.repo.GetBySubjectID(ctx, subjectID)
}
