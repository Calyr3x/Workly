// Работа с данными пользователя: получение, обновление
package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

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

	http.HandleFunc("/updateAvatar", withCORS(handleUpdateAvatar))
	http.HandleFunc("/getUserData", withCORS(handleGetUserData))
	http.HandleFunc("/updateUsername", withCORS(handleUpdateUsername))

	log.Println("Server is running on http://localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

// Middleware для добавления CORS заголовков
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
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

// Получение данных пользователя: юзернейм, аватар, почта
func handleGetUserData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var userData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Avatar   string `json:"avatar"`
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

	err = db.QueryRow("SELECT avatar FROM users WHERE id = $1", userID).Scan(&userData.Avatar)
	if err != nil {
		http.Error(w, "Failed to get avatar", http.StatusInternalServerError)
		return
	}

	// Отправить данные пользователя
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userData)
}
