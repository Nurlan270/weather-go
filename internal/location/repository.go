package location

import (
	"database/sql"
	"fmt"

	"github.com/Nurlan270/weather-go/internal/entity"
)

type DBRepository struct {
	db *sql.DB
}

func NewDBRepository(db *sql.DB) *DBRepository {
	return &DBRepository{db: db}
}

func (r *DBRepository) ListLocationsByUserID(userID int64) ([]*entity.Location, error) {
	const q = "SELECT name, latitude, longitude FROM locations WHERE user_id = $1 ORDER BY id DESC"

	rows, err := r.db.Query(q, userID)
	if err != nil {
		return nil, fmt.Errorf("query lat, lon: %w", err)
	}

	var locationsList []*entity.Location

	for rows.Next() {
		var (
			name     string
			lat, lon float64
		)

		if err = rows.Scan(&name, &lat, &lon); err != nil {
			return nil, fmt.Errorf("get lat, lon: %w", err)
		}

		location := &entity.Location{
			Name: name,
			Coordinates: entity.Coordinates{
				Lat: lat,
				Lon: lon,
			},
		}

		locationsList = append(locationsList, location)
	}

	return locationsList, nil
}

func (r *DBRepository) CreateLocation(userID int64, name string, lat, lon float64) error {
	const q = "INSERT INTO locations (user_id, name, latitude, longitude) VALUES ($1, $2, $3, $4)"

	if _, err := r.db.Exec(q, userID, name, lat, lon); err != nil {
		return fmt.Errorf("insert location: %w", err)
	}

	return nil
}

func (r *DBRepository) DeleteLocation(userID int64, name string) error {
	const q = "DELETE FROM locations WHERE user_id = $1 AND name = $2"

	if _, err := r.db.Exec(q, userID, name); err != nil {
		return fmt.Errorf("delete location: %w", err)
	}

	return nil
}
