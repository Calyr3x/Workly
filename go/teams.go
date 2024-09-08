// Работа с командами
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

	http.HandleFunc("/createTeam", withCORS(handleCreateTeam))
	http.HandleFunc("/getUserAvatar", withCORS(handleGetUserAvatar))

	log.Println("Server is running on http://localhost:8083")
	log.Fatal(http.ListenAndServe(":8083", nil))
}

// Middleware для добавления CORS заголовков
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// Обработчик для создания команды с участниками
func handleCreateTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Name    string   `json:"name"`
		Members []string `json:"members"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID := r.URL.Query().Get("user_id") // Получаем идентификатор пользователя из параметров URL

	// Создать команду в базе данных
	var teamID int
	err := db.QueryRow("INSERT INTO teams (name, owner_id) VALUES ($1, $2) RETURNING id", requestData.Name, userID).Scan(&teamID)
	if err != nil {
		http.Error(w, "Failed to create team", http.StatusInternalServerError)
		return
	}

	// Добавляем участников по их юзернеймам
	for _, username := range requestData.Members {
		var userID uuid.UUID
		err = db.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
		if err != nil {
			http.Error(w, "Failed to find user", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("INSERT INTO team_members (team_id, user_id) VALUES ($1, $2)", teamID, userID)
		if err != nil {
			http.Error(w, "Failed to add member to team", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// Обработчик для получения аватара пользователя по юзернейму
func handleGetUserAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.URL.Query().Get("username") // Получаем юзернейм из параметров URL

	var avatar string
	err := db.QueryRow("SELECT avatar FROM users WHERE username = $1", username).Scan(&avatar)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get avatar", http.StatusInternalServerError)
		}
		return
	}

	// Отправить данные аватара пользователя
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"avatar": avatar,
	})
}
