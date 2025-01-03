package usecase

import (
	"workly/domain"

	"github.com/google/uuid"
)

// TeamRepository интерфейс для работы с данными команд
type TeamRepository interface {
	CreateTeam(name string, ownerID uuid.UUID) (int, error)
	AddMember(teamID int, userID uuid.UUID) error
	RemoveMember(teamID int, userID uuid.UUID) error
	GetUserIDByUsername(username string) (uuid.UUID, error)
	GetTeamsByUserID(userID uuid.UUID) ([]domain.Team, error)
	GetUserAvatar(username string) (string, error)
	IsMemberExists(teamID int, userID uuid.UUID) (bool, error)
}

// TeamUseCase реализует бизнес-логику работы с командами
type TeamUseCase struct {
	repo TeamRepository
}

// NewTeamUseCase создаёт новый экземпляр TeamUseCase
func NewTeamUseCase(repo TeamRepository) *TeamUseCase {
	return &TeamUseCase{repo: repo}
}

func (uc *TeamUseCase) CreateTeam(name string, ownerID uuid.UUID, members []string) (int, error) {
	teamID, err := uc.repo.CreateTeam(name, ownerID)
	if err != nil {
		return 0, err
	}

	// Добавляем создателя команды
	if err := uc.repo.AddMember(teamID, ownerID); err != nil {
		return 0, err
	}

	// Добавляем остальных участников
	for _, username := range members {
		userID, err := uc.repo.GetUserIDByUsername(username)
		if err != nil {
			return 0, err
		}
		if err := uc.repo.AddMember(teamID, userID); err != nil {
			return 0, err
		}
	}

	return teamID, nil
}

func (uc *TeamUseCase) GetTeams(userID uuid.UUID) ([]domain.Team, error) {
	return uc.repo.GetTeamsByUserID(userID)
}

func (uc *TeamUseCase) AddMember(teamID int, username string) error {
	userID, err := uc.repo.GetUserIDByUsername(username)
	if err != nil {
		return err
	}

	exists, err := uc.repo.IsMemberExists(teamID, userID)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrMemberAlreadyExists
	}

	return uc.repo.AddMember(teamID, userID)
}

func (uc *TeamUseCase) RemoveMember(teamID int, username string) error {
	userID, err := uc.repo.GetUserIDByUsername(username)
	if err != nil {
		return err
	}
	return uc.repo.RemoveMember(teamID, userID)
}

func (uc *TeamUseCase) GetUserAvatar(username string) (string, error) {
	return uc.repo.GetUserAvatar(username)
}
