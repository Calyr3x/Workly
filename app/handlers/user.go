package handlers

import (
	"encoding/json"
	"net/http"
	"workly/db"
)

func WithCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func HandleUpdateAvatar(w http.ResponseWriter, r *http.Request) {
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

	userID := r.URL.Query().Get("user_id")
	_, err := db.DB.Exec("UPDATE users SET avatar = $1 WHERE id = $2", profileUpdate.Avatar, userID)
	if err != nil {
		http.Error(w, "Failed to update avatar", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleUpdateUsername(w http.ResponseWriter, r *http.Request) {
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

	userID := r.URL.Query().Get("user_id")
	_, err := db.DB.Exec("UPDATE users SET username = $1 WHERE id = $2", profileUpdate.NewUsername, userID)
	if err != nil {
		http.Error(w, "Failed to update username", http.StatusInternalServerError)
		return
	}
}

func HandleGetUserData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var userData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Avatar   string `json:"avatar"`
	}

	userID := r.URL.Query().Get("user_id")

	err := db.DB.QueryRow("SELECT username FROM users WHERE id = $1", userID).Scan(&userData.Username)
	if err != nil {
		http.Error(w, "Failed to get username", http.StatusInternalServerError)
		return
	}

	err = db.DB.QueryRow("SELECT email FROM users WHERE id = $1", userID).Scan(&userData.Email)
	if err != nil {
		http.Error(w, "Failed to get email", http.StatusInternalServerError)
		return
	}

	err = db.DB.QueryRow("SELECT avatar FROM users WHERE id = $1", userID).Scan(&userData.Avatar)
	if err != nil {
		http.Error(w, "Failed to get avatar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userData)
}
