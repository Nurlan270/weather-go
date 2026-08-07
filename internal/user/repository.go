package user

import (
	"database/sql"
	"fmt"
	"github.com/Nurlan270/weather-go/internal/entity"
	"time"
)

type DBRepository struct {
	db *sql.DB
}

func NewDBRepository(db *sql.DB) *DBRepository {
	return &DBRepository{db: db}
}

func (r *DBRepository) GetUserBySID(sid string) (*entity.User, error) {
	const q = `
		SELECT u.id, u.login, u.password
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.id = $1 AND s.expires_at > now()
	`

	var user entity.User
	err := r.db.QueryRow(q, sid).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *DBRepository) GetUserByLogin(login string) (*entity.User, error) {
	const q = "SELECT id, login, password FROM users WHERE login = $1"

	row := r.db.QueryRow(q, login)

	var user entity.User
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *DBRepository) CreateUser(login, password string) (*entity.User, error) {
	const q = `
		INSERT INTO users (login, password)
		VALUES ($1, $2)
		RETURNING id, login
	`

	row := r.db.QueryRow(q, login, password)

	var user entity.User
	if err := row.Scan(&user.ID, &user.Login); err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return &user, nil
}

func (r *DBRepository) CreateSession(uuid string, userID int64, expiresAt time.Time) (*entity.Session, error) {
	const q = `
		INSERT INTO sessions (id, user_id, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, expires_at
	`

	row := r.db.QueryRow(q, uuid, userID, expiresAt)

	var session entity.Session
	if err := row.Scan(&session.ID, &session.UserID, &session.ExpiresAt); err != nil {
		return nil, fmt.Errorf("failed to insert session: %w", err)
	}

	return &session, nil
}
