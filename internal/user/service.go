package user

import (
	"github.com/Nurlan270/weather-go/internal/entity"
)

type Repository interface {
	GetUserBySID(sid string) (*entity.User, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
