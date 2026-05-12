package models

import "errors"

var (
	ErrNotFound      = errors.New("record not found")
	ErrAlreadyExists = errors.New("record already exists")

	ErrSlotPast        = errors.New("нельзя создать слот в прошлом")
	ErrSlotOverlap     = errors.New("у репетитора уже есть занятие на это время")
	ErrSlotUnavailable = errors.New("слот не существует или уже занят")
	ErrCartEmpty       = errors.New("корзина пуста")
)
