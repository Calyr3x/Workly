package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid" // Импортируем PostgreSQL драйвер
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
	connStr := "user=postgres dbname=Workly sslmode=disable password=calyrexx2003"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	http.HandleFunc("/login", withCORS(handleLogin))
	http.HandleFunc("/register", withCORS(handleRegister))
	http.HandleFunc("/", homePage)                               // Обработчик для корневого пути
	http.HandleFunc("/tasks", withCORS(handleTasks))             // API для получения задач
	http.HandleFunc("/tasks/create", withCORS(handleCreateTask)) // API для создания задач
	http.HandleFunc("/tasks/update", withCORS(handleUpdateTask)) // API для обновления задачи
	http.HandleFunc("/tasks/delete", withCORS(handleDeleteTask)) // API для удаления задачи
	http.HandleFunc("/tasks/", withCORS(handleGetTaskByID))      // Обратите внимание на слэш в пути

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Middleware для добавления CORS заголовков
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// Обработчик главной страницы
func homePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the home page!"))
}

// Обработчик для входа в систему
func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var storedPassword string
	var userID uuid.UUID
	err := db.QueryRow("SELECT id, password FROM users WHERE email = $1", user.Email).Scan(&userID, &storedPassword)
	if err != nil || user.Password != storedPassword {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: userID.String(),
		Path:  "/",
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"user_id": userID.String()})
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.Password)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Обработчик для получения задач пользователя
func handleTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id") // Получаем идентификатор пользователя из параметров URL

	rows, err := db.Query(`
        SELECT t.id, t.name, t.description, t.deadline 
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
		if err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.Deadline); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// Обработчик для создания новой задачи
func handleCreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var task struct {
		Name        string      `json:"name"`
		Description string      `json:"description"`
		Deadline    string      `json:"deadline"` // Принимаем как строку
		CreatorID   uuid.UUID   `json:"creator_id"`
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

	// Вставляем задачу в базу данных и получаем ID новой задачи
	var taskID uuid.UUID
	err = db.QueryRow(`
		INSERT INTO tasks (name, description, deadline, creator_id)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		task.Name, task.Description, deadline, task.CreatorID).Scan(&taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Если список пользователей пустой, добавляем только текущего пользователя
	if len(task.UserIDs) == 0 {
		task.UserIDs = []uuid.UUID{task.CreatorID}
	}

	// Вставляем доступы к задаче для пользователей
	for _, userID := range task.UserIDs {
		_, err := db.Exec(`
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

// Структура задачи
type Task struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Deadline    time.Time   `json:"deadline"`
	CreatorID   uuid.UUID   `json:"creator_id"`
	UserIDs     []uuid.UUID `json:"user_ids"`
}

// Обработчик для редактирования задачи
func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
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

	_, err = db.Exec(`
        UPDATE tasks
        SET name = $1, description = $2, deadline = $3
        WHERE id = $4`,
		task.Name, task.Description, deadline, task.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, userID := range task.UserIDs {
		_, err := db.Exec(`
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
func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(`
        DELETE FROM tasks WHERE id = $1`,
		taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`DELETE FROM task_access WHERE task_id = $1`, taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Обработчик для получения задачи по ID
func handleGetTaskByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Path[len("/tasks/"):] // Извлекаем ID задачи из URL

	var task Task
	err := db.QueryRow(`
        SELECT id, name, description, deadline, creator_id
        FROM tasks
        WHERE id = $1`, taskID).Scan(&task.ID, &task.Name, &task.Description, &task.Deadline, &task.CreatorID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Получаем доступные пользователи для задачи
	rows, err := db.Query(`SELECT user_id FROM task_access WHERE task_id = $1`, taskID)
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
