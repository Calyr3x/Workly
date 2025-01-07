package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"workly/usecase"

	"github.com/google/uuid"
)

type UserHandler struct {
	uc *usecase.UserUseCase
}

func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Avatar string `json:"avatar"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	userID, _ := uuid.Parse(r.URL.Query().Get("user_id"))
	err := h.uc.UpdateAvatar(userID, payload.Avatar)
	if err != nil {
		http.Error(w, "Failed to update avatar", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) GetUserData(w http.ResponseWriter, r *http.Request) {
	userID, _ := uuid.Parse(r.URL.Query().Get("user_id"))
	user, err := h.uc.GetUserData(userID)
	if err != nil {
		http.Error(w, "Failed to get user data", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) UpdateUsername(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	userID, _ := uuid.Parse(r.URL.Query().Get("user_id"))
	err := h.uc.UpdateUsername(userID, payload.Username)
	if err != nil {
		http.Error(w, "Failed to update username", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) GetUserIDs(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Usernames []string `json:"usernames"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	users, err := h.uc.GetUserIDs(payload.Usernames)
	if err != nil {
		http.Error(w, "Failed to get user IDs", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	user, err := h.uc.Login(payload.Email, payload.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"user_id": user.ID.String()})
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	err := h.uc.Register(payload.Email, payload.Password, payload.Username)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
