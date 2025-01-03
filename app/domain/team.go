package domain

import (
	"errors"
	"github.com/google/uuid"
)

// Team представляет команду.
type Team struct {
	ID      int
	Name    string
	OwnerID uuid.UUID
	Members []string
}

// TeamMember представляет участника команды.
type TeamMember struct {
	TeamID int
	UserID uuid.UUID
}

var ErrMemberAlreadyExists = errors.New("member already exists")
