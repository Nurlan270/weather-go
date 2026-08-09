package db

import (
	"database/sql"
	"fmt"

	"github.com/Nurlan270/weather-go/internal/config"

	_ "github.com/lib/pq"
)

func Connect(cfg config.DB) (*sql.DB, error) {
	//	postgres://username:password@host:port/dbname?sslmode=disable
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username, cfg.Password,
		cfg.Host, cfg.Port,
		cfg.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
