package repository

import (
	"database/sql"
	"workly/domain"

	"github.com/google/uuid"
)

type TaskRepositoryImpl struct {
	db *sql.DB
}

// NewTaskRepository создаёт новый экземпляр TaskRepository.
func NewTaskRepository(db *sql.DB) *TaskRepositoryImpl {
	return &TaskRepositoryImpl{db: db}
}

func (r *TaskRepositoryImpl) GetTasksByUserID(userID uuid.UUID) ([]domain.Task, error) {
	rows, err := r.db.Query(`
		SELECT t.id, t.name, t.description, t.deadline, t.created_at, t.creator_id, t.status
		FROM tasks t
		INNER JOIN task_access ta ON t.id = ta.task_id
		WHERE ta.user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.Deadline, &task.CreatedAt, &task.CreatorID, &task.Status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepositoryImpl) CreateTask(task domain.Task) (uuid.UUID, error) {
	err := r.db.QueryRow(`
		INSERT INTO tasks (id, name, description, deadline, creator_id, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		task.ID, task.Name, task.Description, task.Deadline, task.CreatorID, task.Status, task.CreatedAt).Scan(&task.ID)
	return task.ID, err
}

func (r *TaskRepositoryImpl) UpdateTask(task domain.Task) error {
	_, err := r.db.Exec(`
		UPDATE tasks
		SET name = $1, description = $2, deadline = $3, status = $4
		WHERE id = $5`,
		task.Name, task.Description, task.Deadline, task.Status, task.ID)
	return err
}

func (r *TaskRepositoryImpl) DeleteTask(taskID uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM tasks WHERE id = $1`, taskID)
	return err
}

func (r *TaskRepositoryImpl) GetTaskByID(taskID uuid.UUID) (*domain.Task, error) {
	var task domain.Task
	err := r.db.QueryRow(`
		SELECT id, name, description, deadline, created_at, creator_id, status
		FROM tasks
		WHERE id = $1`, taskID).Scan(&task.ID, &task.Name, &task.Description, &task.Deadline, &task.CreatedAt, &task.CreatorID, &task.Status)
	return &task, err
}

func (r *TaskRepositoryImpl) CreateTaskAccess(access domain.TaskAccess) error {
	_, err := r.db.Exec(`
		INSERT INTO task_access (task_id, user_id)
		VALUES ($1, $2)`,
		access.TaskID, access.UserID)
	return err
}
