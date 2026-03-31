package models

import "time"

// Session представляет активную сессию пользователя
type Session struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
}
