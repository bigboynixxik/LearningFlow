package models

type User struct {
	Id         string
	Name       string
	Email      string
	Password   string
	Login      string
	SubjectsId []int64
}
