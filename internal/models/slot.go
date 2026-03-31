package models

import "time"

// Slot - временной интервал репетитора. Выступает также в роли "товара" для корзины.
type Slot struct {
	ID        string
	TutorID   string
	StudentID *string
	StartTime time.Time
	EndTime   time.Time
	Status    string // "free", "in_cart", "booked"
}
