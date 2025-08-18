package domain

import "time"

// User описывает сущность пользователя
type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username" validate:"required,min=3,max=100"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
