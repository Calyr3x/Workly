package repository

import (
	"database/sql"
	"workly/domain"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type TeamRepositoryImpl struct {
	db *sql.DB
}

// NewTeamRepository создаёт новый экземпляр TeamRepositoryImpl.
func NewTeamRepository(db *sql.DB) *TeamRepositoryImpl {
	return &TeamRepositoryImpl{db: db}
}

func (r *TeamRepositoryImpl) CreateTeam(name string, ownerID uuid.UUID) (int, error) {
	var teamID int
	err := r.db.QueryRow("INSERT INTO teams (name, owner_id) VALUES ($1, $2) RETURNING id", name, ownerID).Scan(&teamID)
	return teamID, err
}

func (r *TeamRepositoryImpl) AddMember(teamID int, userID uuid.UUID) error {
	_, err := r.db.Exec("INSERT INTO team_members (team_id, user_id) VALUES ($1, $2)", teamID, userID)
	return err
}

func (r *TeamRepositoryImpl) RemoveMember(teamID int, userID uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM team_members WHERE team_id = $1 AND user_id = $2", teamID, userID)
	return err
}

func (r *TeamRepositoryImpl) GetUserIDByUsername(username string) (uuid.UUID, error) {
	var userID uuid.UUID
	err := r.db.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	return userID, err
}

func (r *TeamRepositoryImpl) GetTeamsByUserID(userID uuid.UUID) ([]domain.Team, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT t.id, t.name, t.owner_id, array_agg(DISTINCT u.username)
		FROM teams t
		JOIN team_members tm ON t.id = tm.team_id
		JOIN users u ON tm.user_id = u.id
		WHERE t.id IN (
			SELECT t1.id
			FROM teams t1
			JOIN team_members tm1 ON t1.id = tm1.team_id
			WHERE t1.owner_id = $1 OR tm1.user_id = $1
		)
		GROUP BY t.id, t.name, t.owner_id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []domain.Team
	for rows.Next() {
		var team domain.Team
		err := rows.Scan(&team.ID, &team.Name, &team.OwnerID, pq.Array(&team.Members))
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}

func (r *TeamRepositoryImpl) GetUserAvatar(username string) (string, error) {
	var avatar string
	err := r.db.QueryRow("SELECT avatar FROM users WHERE username = $1", username).Scan(&avatar)
	return avatar, err
}

func (r *TeamRepositoryImpl) IsMemberExists(teamID int, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM team_members WHERE team_id = $1 AND user_id = $2)", teamID, userID).Scan(&exists)
	return exists, err
}
