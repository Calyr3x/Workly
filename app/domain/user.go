package domain

import (
	"errors"
	"github.com/google/uuid"
)

// User структура пользователя
type User struct {
	ID       uuid.UUID
	Email    string
	Password string
	Username string
	Avatar   string
}

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserNotFound       = errors.New("user not found")
)
