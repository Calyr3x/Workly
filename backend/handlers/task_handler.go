package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"workly/domain"
	"workly/usecase"

	"github.com/google/uuid"
)

type TaskHandler struct {
	uc *usecase.TaskUseCase
}

func NewTaskHandler(uc *usecase.TaskUseCase) *TaskHandler {
	return &TaskHandler{uc: uc}
}

// GetTasks получает все доступные для пользователя задачи
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	userID, _ := uuid.Parse(r.URL.Query().Get("user_id"))
	tasks, err := h.uc.GetTasks(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// CreateTask создает задачу
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Deadline    string    `json:"deadline"`
		CreatorID   uuid.UUID `json:"creator_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	deadline, err := time.Parse(time.RFC3339, payload.Deadline)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	taskID, err := h.uc.CreateTask(payload.Name, payload.Description, deadline, payload.CreatorID)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	if err = h.uc.CreateTaskAccess(taskID, payload.CreatorID); err != nil {
		http.Error(w, "Failed to create task access", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"task_id": taskID.String()})
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// UpdateTask обновляет задачу
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Deadline    string    `json:"deadline"`
		Status      string    `json:"taskStatus"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	deadline, err := time.Parse(time.RFC3339, payload.Deadline)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	task := domain.Task{
		ID:          payload.ID,
		Name:        payload.Name,
		Description: payload.Description,
		Deadline:    deadline,
		Status:      payload.Status,
	}

	if err := h.uc.UpdateTask(task); err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"task_id": task.ID.String()})
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteTask удаляет задачу
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := h.uc.DeleteTask(taskID); err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetTaskByID возвращает задачу по ID
func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	taskID, err := uuid.Parse(r.URL.Path[len("/tasks/"):])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.uc.GetTaskByID(taskID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve task", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// CreateTaskAccess создаёт доступ к задаче
func (h *TaskHandler) CreateTaskAccess(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		TaskID uuid.UUID `json:"task_id"`
		UserID uuid.UUID `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.uc.CreateTaskAccess(payload.TaskID, payload.UserID); err != nil {
		http.Error(w, "Failed to create task access", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
