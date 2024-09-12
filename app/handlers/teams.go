package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"workly/db"

	"github.com/google/uuid"
)

// Обработчик для создания команды с участниками
func HandleCreateTeam(w http.ResponseWriter, r *http.Request) {
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
	err := db.DB.QueryRow("INSERT INTO teams (name, owner_id) VALUES ($1, $2) RETURNING id", requestData.Name, userID).Scan(&teamID)
	if err != nil {
		http.Error(w, "Failed to create team", http.StatusInternalServerError)
		return
	}

	// Добавляем участников по их юзернеймам
	for _, username := range requestData.Members {
		var memberID uuid.UUID
		err = db.DB.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&memberID)
		if err != nil {
			http.Error(w, "Failed to find user", http.StatusBadRequest)
			return
		}

		_, err = db.DB.Exec("INSERT INTO team_members (team_id, user_id) VALUES ($1, $2)", teamID, memberID)
		if err != nil {
			http.Error(w, "Failed to add member to team", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// Обработчик для получения аватара пользователя по юзернейму
func HandleGetUserAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.URL.Query().Get("username") // Получаем юзернейм из параметров URL

	var avatar string
	err := db.DB.QueryRow("SELECT avatar FROM users WHERE username = $1", username).Scan(&avatar)
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
