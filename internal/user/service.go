package user

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Nurlan270/weather-go/internal/config"
	"github.com/Nurlan270/weather-go/internal/entity"
	"github.com/Nurlan270/weather-go/internal/rest/request"
	guuid "github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserRepository interface {
	GetUserBySID(sid string) (*entity.User, error)
	GetUserByLogin(login string) (*entity.User, error)
	CreateUser(login, password string) (*entity.User, error)
	CreateSession(uuid string, userID int64, expiresAt time.Time) (*entity.Session, error)
}

type Service struct {
	repo    UserRepository
	sessCfg config.Session
}

func NewService(repo UserRepository, sessCfg config.Session) *Service {
	return &Service{
		repo:    repo,
		sessCfg: sessCfg,
	}
}

func (s *Service) GetUserBySID(sid string) (*entity.User, error) {
	u, err := s.repo.GetUserBySID(sid)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrSessionExpired
	}

	return u, err
}

func (s *Service) GetUserByLogin(login string) (*entity.User, error) {
	u, err := s.repo.GetUserByLogin(login)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	return u, err
}

func (s *Service) RegisterUser(user request.RegisterUser) (*entity.Session, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	expiresAt := time.Now().Add(s.sessCfg.ExpiresIn).UTC()

	//	Create User
	u, err := s.repo.CreateUser(user.Login, string(hashedPass))
	if err != nil {
		uniqErr := pq.As(err, pqerror.UniqueViolation)
		if uniqErr != nil {
			return nil, ErrUserAlreadyExists
		}

		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	uuid, err := s.generateUUID()
	if err != nil {
		return nil, err
	}

	//	Create Session
	session, err := s.repo.CreateSession(uuid, u.ID, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func (s *Service) LoginUser(user request.LoginUser) (*entity.Session, error) {
	u, err := s.repo.GetUserByLogin(user.Login)

	//	Check if user exists
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	//	Check if password is valid
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)); err != nil {
		return nil, ErrInvalidPassword
	}

	//	All challenges passed -> Login User (create new session)
	uuid, err := s.generateUUID()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(s.sessCfg.ExpiresIn).UTC()

	session, err := s.repo.CreateSession(uuid, u.ID, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func (s *Service) generateUUID() (string, error) {
	uuid, err := guuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate uuid: %w", err)
	}

	return uuid.String(), err
}
