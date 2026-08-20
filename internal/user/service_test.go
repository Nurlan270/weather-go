package user

import (
	"context"
	"database/sql"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"golang.org/x/crypto/bcrypt"

	"github.com/Nurlan270/weather-go/internal/config"
	database "github.com/Nurlan270/weather-go/internal/db"
	"github.com/Nurlan270/weather-go/internal/entity"
	"github.com/Nurlan270/weather-go/internal/rest/request"
)

func TestUserService_RegisterUser(t *testing.T) {
	//	Arrange
	var (
		userLogin    = "foo@bar.com"
		userPassword = "secret-123"

		sessConf = config.Session{
			Name:      "test_session",
			ExpiresIn: 5 * time.Minute,
		}

		registerUser = request.RegisterUser{
			Login:                userLogin,
			Password:             userPassword,
			PasswordConfirmation: userPassword,
		}

		mockUserRepo = NewMockUserRepository(t)
	)

	//	Tests
	t.Run("it registers valid user", func(t *testing.T) {
		//	Arrange
		var (
			expectedUser = &entity.User{
				ID:       1,
				Login:    userLogin,
				Password: userPassword,
			}

			expectedSession = &entity.Session{
				ID:        "some-uuid",
				UserID:    expectedUser.ID,
				ExpiresAt: time.Now().Add(sessConf.ExpiresIn).UTC(),
			}

			userSvc = NewService(mockUserRepo, sessConf)
		)

		//	Mocks
		mockUserRepo.
			On("CreateUser", userLogin, mock.MatchedBy(func(hashed string) bool {
				return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)) == nil
			})).
			Once().
			Return(expectedUser, nil)

		mockUserRepo.
			On(
				"CreateSession",
				mock.AnythingOfType("string"),
				expectedUser.ID,
				mock.AnythingOfType("time.Time"),
			).
			Once().
			Return(expectedSession, nil)

		//	Act
		actualSession, err := userSvc.RegisterUser(registerUser)

		require.NoError(t, err, "Should be able to register user")
		require.Equal(t, expectedSession, actualSession, "Should return expected session")
	})

	t.Run("it does not register user if it already exists", func(t *testing.T) {
		//	Arrange
		db := initDB(t)
		dbUserRepo := NewDBRepository(db)
		userSvc := NewService(dbUserRepo, sessConf)

		//	Act
		_, err := userSvc.RegisterUser(registerUser)

		require.NoError(t, err, "Should be able to register user")

		_, err = userSvc.RegisterUser(registerUser)

		require.NotNil(t, err, "Should be a valid error")
		require.ErrorIs(t, err, ErrUserAlreadyExists, "Should be unable to register user with same login")
	})
}

func initDB(t *testing.T) *sql.DB {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//	Arrange
	const PostgresImg = "postgres:18"

	var (
		dbName = "test_db"
		dbUser = "user"
		dbPass = "secret-123"
	)

	//	Test Container
	pgc, err := postgres.Run(ctx, PostgresImg,
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	t.Cleanup(func() {
		if err = testcontainers.TerminateContainer(pgc); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	//	Get Host/Port pair from container
	dbHostPort, err := pgc.Endpoint(ctx, "")
	if err != nil {
		t.Fatalf("failed to retrieve postgres host/port pair: %s", err)
	}

	dbHost, dbPort, err := net.SplitHostPort(dbHostPort)
	if err != nil {
		t.Fatalf("failed to parse postgres host and port: %s", err)
	}

	dbConf := config.DB{
		Host:     dbHost,
		Port:     dbPort,
		Name:     dbName,
		Username: dbUser,
		Password: dbPass,
	}

	db, err := database.Connect(dbConf)
	if err != nil {
		t.Fatalf("failed to connect to database: %s", err)
	}

	//	Setup Goose
	if err = goose.SetDialect("postgres"); err != nil {
		t.Fatalf("failed to set goose dialect: %s", err)
	}

	//	Run all migrations
	if err = goose.Up(db, filepath.Join("..", "db", "migrations")); err != nil {
		t.Fatalf("failed to run migrations: %s", err)
	}

	return db
}
