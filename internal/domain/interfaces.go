package domain

import (
	"context"
	"learningflow/internal/models"
)

// UserRepository управляет учетными записями
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

// TutorRepository отдает данные для витрины (Лаба №1)
type TutorRepository interface {
	GetAll(ctx context.Context) ([]models.Tutor, error)
	GetByID(ctx context.Context, id string) (*models.Tutor, error)
	GetBySubjectID(ctx context.Context, subjectID int64) ([]models.Tutor, error)
}

// SubjectRepository для вывода меню предметов
type SubjectRepository interface {
	GetAll(ctx context.Context) ([]models.Subject, error)
}

// SlotRepository управляет расписанием и корзиной
type SlotRepository interface {
	Create(ctx context.Context, slot *models.Slot) error

	GetFreeSlots(ctx context.Context, tutorID string) ([]models.Slot, error)

	AddToCart(ctx context.Context, slotID string, studentID string) error

	GetCart(ctx context.Context, studentID string) ([]models.Slot, error)

	RemoveFromCart(ctx context.Context, slotID string, studentID string) error

	CheckoutCart(ctx context.Context, studentID string) error
}
