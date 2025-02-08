package domain

import (
	"github.com/google/uuid"
	"time"
)

// Task структура задачи.
type Task struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	Deadline    time.Time
	CreatorID   uuid.UUID
	Status      string
	Categories  []Category
}

// Category структура категории задачи
type Category struct {
	Name  string
	Color string
}

// TaskAccess структура доступа к задаче.
type TaskAccess struct {
	TaskID uuid.UUID
	UserID uuid.UUID
}
