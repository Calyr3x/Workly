package usecase

import (
	"time"
	"workly/domain"

	"github.com/google/uuid"
)

// TaskRepository описывает доступ к данным задач.
type TaskRepository interface {
	GetTasksByUserID(userID uuid.UUID) ([]domain.Task, error)
	CreateTask(task domain.Task) (uuid.UUID, error)
	UpdateTask(task domain.Task) error
	DeleteTask(taskID uuid.UUID) error
	GetTaskByID(taskID uuid.UUID) (*domain.Task, error)
	CreateTaskAccess(access domain.TaskAccess) error
}

// TaskUseCase реализует бизнес-логику задач.
type TaskUseCase struct {
	repo TaskRepository
}

// NewTaskUseCase создаёт новый экземпляр TaskUseCase.
func NewTaskUseCase(repo TaskRepository) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}

func (uc *TaskUseCase) GetTasks(userID uuid.UUID) ([]domain.Task, error) {
	return uc.repo.GetTasksByUserID(userID)
}

func (uc *TaskUseCase) CreateTask(name, description string, deadline time.Time, creatorID uuid.UUID) (uuid.UUID, error) {
	task := domain.Task{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		Deadline:    deadline,
		CreatorID:   creatorID,
		Status:      "new",
	}
	return uc.repo.CreateTask(task)
}

func (uc *TaskUseCase) UpdateTask(task domain.Task) error {
	return uc.repo.UpdateTask(task)
}

func (uc *TaskUseCase) DeleteTask(taskID uuid.UUID) error {
	return uc.repo.DeleteTask(taskID)
}

func (uc *TaskUseCase) GetTaskByID(taskID uuid.UUID) (*domain.Task, error) {
	return uc.repo.GetTaskByID(taskID)
}

func (uc *TaskUseCase) CreateTaskAccess(taskID, userID uuid.UUID) error {
	access := domain.TaskAccess{TaskID: taskID, UserID: userID}
	return uc.repo.CreateTaskAccess(access)
}
