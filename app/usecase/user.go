package usecase

import (
	"github.com/google/uuid"
	"workly/domain"
)

// UserRepository описывает доступ к данным пользователей.
type UserRepository interface {
	UpdateAvatar(userID uuid.UUID, avatar string) error
	UpdateUsername(userID uuid.UUID, username string) error
	GetUserByID(userID uuid.UUID) (*domain.User, error)
	GetUserIDsByUsernames(usernames []string) ([]domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	Create(user domain.User) error
}

// UserUseCase описывает бизнес-логику работы с пользователями.
type UserUseCase struct {
	repo UserRepository
}

// NewUserUseCase создает новый экземпляр UserUseCase.
func NewUserUseCase(repo UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) UpdateAvatar(userID uuid.UUID, avatar string) error {
	return uc.repo.UpdateAvatar(userID, avatar)
}

func (uc *UserUseCase) UpdateUsername(userID uuid.UUID, username string) error {
	return uc.repo.UpdateUsername(userID, username)
}

func (uc *UserUseCase) GetUserData(userID uuid.UUID) (*domain.User, error) {
	return uc.repo.GetUserByID(userID)
}

func (uc *UserUseCase) GetUserIDs(usernames []string) ([]domain.User, error) {
	return uc.repo.GetUserIDsByUsernames(usernames)
}

func (uc *UserUseCase) Login(email, password string) (*domain.User, error) {
	user, err := uc.repo.FindByEmail(email)
	if err != nil || user.Password != password {
		return nil, domain.ErrInvalidCredentials
	}
	return user, nil
}

func (uc *UserUseCase) Register(email, password, username string) error {
	user := domain.User{
		ID:       uuid.New(),
		Email:    email,
		Password: password,
		Username: username,
	}
	return uc.repo.Create(user)
}
