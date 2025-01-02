package repository

import (
	"database/sql"
	"errors"
	"workly/domain"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) UpdateAvatar(userID uuid.UUID, avatar string) error {
	_, err := r.db.Exec("UPDATE users SET avatar = $1 WHERE id = $2", avatar, userID)
	return err
}

func (r *UserRepositoryImpl) UpdateUsername(userID uuid.UUID, username string) error {
	_, err := r.db.Exec("UPDATE users SET username = $1 WHERE id = $2", username, userID)
	return err
}

func (r *UserRepositoryImpl) GetUserByID(userID uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, email, username, avatar FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Email, &user.Username, &user.Avatar)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	return &user, err
}

func (r *UserRepositoryImpl) GetUserIDsByUsernames(usernames []string) ([]domain.User, error) {
	rows, err := r.db.Query("SELECT id, username FROM users WHERE username = ANY($1)", pq.Array(usernames))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, email, password FROM users WHERE email = $1", email).Scan(&user.ID, &user.Email, &user.Password)
	return &user, err
}

func (r *UserRepositoryImpl) Create(user domain.User) error {
	_, err := r.db.Exec("INSERT INTO users (id, email, password, username, avatar) VALUES ($1, $2, $3, $4, $5)", user.ID, user.Email, user.Password, user.Username, user.Avatar)
	return err
}
