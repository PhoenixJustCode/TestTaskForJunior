package domain

import (
	"time"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type User struct {
    ID           int       `json:"id" db:"id"`
    Username     string    `json:"username" db:"username" validate:"required,min=2,max=99"`
    Email        string    `json:"email" db:"email" validate:"required,email"`
    PasswordHash string    `json:"-" db:"password_hash" validate:"required,gte=6"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type SignUpInput struct {
    Name     string `json:"username" validate:"required,min=2,max=99"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,gte=6"`
}

func (i SignUpInput) Validate() error {
    return validate.Struct(i)
}

type SignInInput struct {
    Name     string `json:"username" validate:"required,min=2"`
    Password string `json:"password" validate:"required,gte=6"`
}

func (i SignInInput) Validate() error {
    return validate.Struct(i)
}
