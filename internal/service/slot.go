package service

import (
	"context"
	"learningflow/internal/models"
)

type SlotRepository interface {
	Create(ctx context.Context, slot *models.Slot) error

	GetFreeSlots(ctx context.Context, tutorID string) ([]models.Slot, error)

	AddToCart(ctx context.Context, slotID string, studentID string) error

	GetCart(ctx context.Context, studentID string) ([]models.Slot, error)

	RemoveFromCart(ctx context.Context, slotID string, studentID string) error

	CheckoutCart(ctx context.Context, studentID string) error
}
