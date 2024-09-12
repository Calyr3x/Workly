package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"workly/db"

	"github.com/google/uuid"
)

// Task структура задачи
type Task struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	Deadline    time.Time   `json:"deadline"`
	CreatorID   uuid.UUID   `json:"creator_id"`
	UserIDs     []uuid.UUID `json:"user_ids"`
}

// HandleTasks обрабатывает запросы на получение задач
func HandleTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	rows, err := db.DB.Query(`
	SELECT t.id, t.name, t.description, t.deadline, t.created_at 
	FROM tasks t 
	INNER JOIN task_access ta ON t.id = ta.task_id 
	WHERE ta.user_id = $1`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.Deadline, &task.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// HandleCreateTask обрабатывает создание новой задачи
func HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var task struct {
		Name        string      `json:"name"`
		Description string      `json:"description"`
		Deadline    string      `json:"deadline"`
		CreatorID   uuid.UUID   `json:"creator_id"`
		UserIDs     []uuid.UUID `json:"user_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deadline, err := time.Parse(time.RFC3339, task.Deadline)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	var taskID uuid.UUID
	err = db.DB.QueryRow(`
		INSERT INTO tasks (name, description, deadline, creator_id)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		task.Name, task.Description, deadline, task.CreatorID).Scan(&taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(task.UserIDs) == 0 {
		task.UserIDs = []uuid.UUID{task.CreatorID}
	}

	for _, userID := range task.UserIDs {
		_, err := db.DB.Exec(`
			INSERT INTO task_access (task_id, user_id) VALUES ($1, $2)`,
			taskID, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"task_id": taskID.String()})
}

// Обработчик для редактирования задачи
func HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var task struct {
		ID          uuid.UUID   `json:"id"`
		Name        string      `json:"name"`
		Description string      `json:"description"`
		Deadline    string      `json:"deadline"`
		UserIDs     []uuid.UUID `json:"user_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Преобразуем строку даты в time.Time
	deadline, err := time.Parse(time.RFC3339, task.Deadline)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(`
        UPDATE tasks
        SET name = $1, description = $2, deadline = $3
        WHERE id = $4`,
		task.Name, task.Description, deadline, task.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, userID := range task.UserIDs {
		_, err := db.DB.Exec(`
            INSERT INTO task_access (task_id, user_id) VALUES ($1, $2)`,
			task.ID, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// Обработчик для удаления задачи
func HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	_, err := db.DB.Exec(`
        DELETE FROM tasks WHERE id = $1`,
		taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec(`DELETE FROM task_access WHERE task_id = $1`, taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Обработчик для получения задачи по ID
func HandleGetTaskByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Path[len("/tasks/"):] // Извлекаем ID задачи из URL

	var task Task
	err := db.DB.QueryRow(`
	SELECT id, name, description, deadline, created_at, creator_id
	FROM tasks
	WHERE id = $1`, taskID).Scan(&task.ID, &task.Name, &task.Description, &task.Deadline, &task.CreatedAt, &task.CreatorID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Получаем доступные пользователи для задачи
	rows, err := db.DB.Query(`SELECT user_id FROM task_access WHERE task_id = $1`, taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var userIDs []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userIDs = append(userIDs, userID)
	}
	task.UserIDs = userIDs

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}
