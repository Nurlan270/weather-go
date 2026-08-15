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

func (r *DBRepository) ListLocationIdsByUserID(userID int64) ([]int64, error) {
	const q = `
		SELECT id
		FROM locations
		WHERE user_id = $1
		ORDER BY id DESC
	`

	rows, err := r.db.Query(q, userID)
	if err != nil {
		return nil, fmt.Errorf("query location id: %w", err)
	}

	var ids []int64

	for rows.Next() {
		var id int64

		if err = rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("get location id: %w", err)
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (r *DBRepository) GetLocationByID(id int64) (entity.Location, error) {
	const q = `
		SELECT id, name, latitude, longitude
		FROM locations
		WHERE id = $1
		ORDER BY id DESC
	`

	var location entity.Location
	if err := r.db.QueryRow(q, id).Scan(
		&location.ID,
		&location.Name,
		&location.Coordinates.Lat,
		&location.Coordinates.Lon,
	); err != nil {
		return entity.Location{}, fmt.Errorf("query location: %w", err)
	}

	return location, nil
}

func (r *DBRepository) CreateLocation(userID int64, name string, lat, lon float64) error {
	const q = "INSERT INTO locations (user_id, name, latitude, longitude) VALUES ($1, $2, $3, $4)"

	if _, err := r.db.Exec(q, userID, name, lat, lon); err != nil {
		return fmt.Errorf("insert location: %w", err)
	}

	return nil
}

func (r *DBRepository) DeleteLocation(id, userID int64) error {
	const q = "DELETE FROM locations WHERE id = $1 AND user_id = $2"

	if _, err := r.db.Exec(q, id, userID); err != nil {
		return fmt.Errorf("delete location: %w", err)
	}

	return nil
}
