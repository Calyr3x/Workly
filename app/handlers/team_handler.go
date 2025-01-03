package handlers

import (
	"encoding/json"
	"net/http"
	"workly/usecase"

	"github.com/google/uuid"
)

type TeamHandler struct {
	uc *usecase.TeamUseCase
}

func NewTeamHandler(uc *usecase.TeamUseCase) *TeamHandler {
	return &TeamHandler{uc: uc}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name    string   `json:"name"`
		Members []string `json:"members"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, _ := uuid.Parse(r.URL.Query().Get("user_id"))
	teamID, err := h.uc.CreateTeam(payload.Name, userID, payload.Members)
	if err != nil {
		http.Error(w, "Failed to create team", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"team_id": teamID})
}

func (h *TeamHandler) GetTeams(w http.ResponseWriter, r *http.Request) {
	userID, _ := uuid.Parse(r.URL.Query().Get("user_id"))
	teams, err := h.uc.GetTeams(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve teams", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(teams)
}

func (h *TeamHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		TeamID int      `json:"team_id"`
		Member []string `json:"Member"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Обрабатываем каждого участника
	for _, member := range payload.Member {
		if err := h.uc.AddMember(payload.TeamID, member); err != nil {
			http.Error(w, "Failed to add member", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		TeamID int    `json:"team_id"`
		Member string `json:"member"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.uc.RemoveMember(payload.TeamID, payload.Member); err != nil {
		http.Error(w, "Failed to remove member", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TeamHandler) GetUserAvatar(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	avatar, err := h.uc.GetUserAvatar(username)
	if err != nil {
		http.Error(w, "Failed to retrieve avatar", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"avatar": avatar})
}
