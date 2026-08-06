package user

import (
	"database/sql"
	"github.com/Nurlan270/weather-go/internal/entity"
)

type DBRepository struct {
	db *sql.DB
}

func NewDBRepository(db *sql.DB) *DBRepository {
	return &DBRepository{db: db}
}

func (r *DBRepository) GetUserBySID(sid string) (*entity.User, error) {
	return nil, nil
}
