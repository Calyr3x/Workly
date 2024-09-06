package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
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
	http.HandleFunc("/updateAvatar", withCORS(handleUpdateAvatar))
	http.HandleFunc("/getCurrentAvatar", withCORS(handleGetCurrentAvatar))
	http.HandleFunc("/getUserData", withCORS(handleGetUserData))
	http.HandleFunc("/updateUsername", withCORS(handleUpdateUsername))

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

// Обработчик входа в систему
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"user_id": userID.String()})
}

// Обработка регситрации пользователя
func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json: "username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO users (email, password, username) VALUES ($1, $2, $3)", user.Email, user.Password, user.Username)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Обработчик для обновления аватара пользователя
func handleUpdateAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var profileUpdate struct {
		Avatar string `json:"avatar"`
	}
	if err := json.NewDecoder(r.Body).Decode(&profileUpdate); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID := r.URL.Query().Get("user_id") // Получаем идентификатор пользователя из параметров URL
	// Обновить профиль в базе данных
	_, err := db.Exec("UPDATE users SET avatar = $1 WHERE id = $2", profileUpdate.Avatar, userID)
	if err != nil {
		http.Error(w, "Failed to update avatar", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Обработчик для обновления имени пользователя
func handleUpdateUsername(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var profileUpdate struct {
		NewUsername string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&profileUpdate); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID := r.URL.Query().Get("user_id") // Получаем идентификатор пользователя из параметров URL

	_, err := db.Exec("UPDATE users SET username = $1 WHERE id = $2", profileUpdate.NewUsername, userID)
	if err != nil {
		http.Error(w, "Failed to update username", http.StatusInternalServerError)
		return
	}
}

// Обработчик для получения текущего аватара пользователя
func handleGetCurrentAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id") // Получаем идентификатор пользователя из параметров URL
	var avatar string
	err := db.QueryRow("SELECT avatar FROM users WHERE id = $1", userID).Scan(&avatar)
	if err != nil {
		http.Error(w, "Failed to get avatar", http.StatusInternalServerError)
		return
	}

	// Отправить текущий аватар пользователю
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"avatar": avatar})
}

// Получение юзернейма
func handleGetUserData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var userData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	userID := r.URL.Query().Get("user_id") // Получаем идентификатор пользователя из параметров URL

	err := db.QueryRow("SELECT username FROM users WHERE id = $1", userID).Scan(&userData.Username)
	if err != nil {
		http.Error(w, "Failed to get username", http.StatusInternalServerError)
		return
	}

	err = db.QueryRow("SELECT email FROM users WHERE id = $1", userID).Scan(&userData.Email)
	if err != nil {
		http.Error(w, "Failed to get email", http.StatusInternalServerError)
		return
	}
	// Отправить текущий юзернейм пользователю
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userData)
}
