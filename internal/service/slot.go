package service

import (
	"context"
	"fmt"
	"learningflow/internal/models"
	"time"
)

type SlotRepository interface {
	Create(ctx context.Context, slot *models.Slot) error
	GetFreeSlots(ctx context.Context, tutorID string) ([]models.Slot, error)
	AddToCart(ctx context.Context, slotID string, studentID string) error
	GetCart(ctx context.Context, studentID string) ([]models.Slot, error)
	RemoveFromCart(ctx context.Context, slotID string, studentID string) error
	CheckoutCart(ctx context.Context, studentID string) error
	GetCartCount(ctx context.Context, studentID string) (int, error)
}

type SlotService struct {
	repo SlotRepository
}

func NewSlotService(repo SlotRepository) *SlotService {
	return &SlotService{repo: repo}
}

func (s *SlotService) CreateAvailability(ctx context.Context, tutorID string, startTime time.Time) error {
	if startTime.Before(time.Now()) {
		return models.ErrSlotPast
	}

	endTime := startTime.Add(time.Hour)

	slot := &models.Slot{
		TutorID:   tutorID,
		StartTime: startTime,
		EndTime:   endTime,
	}

	return s.repo.Create(ctx, slot)
}

func (s *SlotService) GetFreeSlots(ctx context.Context, tutorID string) ([]models.Slot, error) {
	return s.repo.GetFreeSlots(ctx, tutorID)
}

func (s *SlotService) AddToCart(ctx context.Context, slotID string, studentID string) error {
	return s.repo.AddToCart(ctx, slotID, studentID)
}

func (s *SlotService) GetCart(ctx context.Context, studentID string) ([]models.Slot, error) {
	return s.repo.GetCart(ctx, studentID)
}

func (s *SlotService) RemoveFromCart(ctx context.Context, slotID string, studentID string) error {
	return s.repo.RemoveFromCart(ctx, slotID, studentID)
}

func (s *SlotService) CheckoutCart(ctx context.Context, studentID string) error {
	count, err := s.repo.GetCartCount(ctx, studentID)
	if err != nil {
		return fmt.Errorf("service.slot.CheckoutCart count: %w", err)
	}
	if count == 0 {
		return models.ErrCartEmpty
	}

	return s.repo.CheckoutCart(ctx, studentID)
}

func (s *SlotService) GetCartCount(ctx context.Context, studentID string) (int, error) {
	return s.repo.GetCartCount(ctx, studentID)
}
