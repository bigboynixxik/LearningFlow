package models

type Role string

const (
	RoleStudent Role = "student"
	RoleTutor   Role = "tutor"
)

// User - базовая сущность пользователя (для авторизации)
type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         Role
}
