package handlers

import (
	"encoding/json"
	"net/http"

	"workly/db"

	"github.com/google/uuid"
)

// HandleLogin обрабатывает вход пользователя
func HandleLogin(w http.ResponseWriter, r *http.Request) {
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
	err := db.DB.QueryRow("SELECT id, password FROM users WHERE email = $1", user.Email).Scan(&userID, &storedPassword)
	if err != nil || user.Password != storedPassword {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"user_id": userID.String()})
}

// HandleRegister обрабатывает регистрацию пользователя
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err := db.DB.Exec("INSERT INTO users (email, password, username) VALUES ($1, $2, $3)", user.Email, user.Password, user.Username)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
