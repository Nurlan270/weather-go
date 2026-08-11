package location

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"

	"github.com/Nurlan270/weather-go/internal/entity"
	"github.com/Nurlan270/weather-go/internal/rest/openweather/response"
	"github.com/Nurlan270/weather-go/internal/rest/request"
)

type LocationRepository interface {
	ListLocationsByUserID(userID int64) ([]*entity.Location, error)
	CreateLocation(userID int64, name string, lat, lon float64) error
	DeleteLocation(userID int64, name string) error
}

type LocationClient interface {
	GetByCityName(query string) (*http.Response, error)
	GetByCityNameAndCoordinates(name, lat, lon string) (*http.Response, error)
}

type Service struct {
	repo   LocationRepository
	client LocationClient
}

func NewService(repo LocationRepository, client LocationClient) *Service {
	return &Service{
		repo:   repo,
		client: client,
	}
}

func (s *Service) SearchLocation(request request.SearchLocation) (*response.Location, error) {
	resp, err := s.client.GetByCityName(request.Query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNoResults
	}

	var location response.Location
	if err = json.NewDecoder(resp.Body).Decode(&location); err != nil {
		return nil, fmt.Errorf("failed to decode location: %w", err)
	}

	//	Round Temperature
	location.Main.Temp = math.Round(location.Main.Temp)

	return &location, err
}

func (s *Service) AddLocation(userID int64, request request.AddLocation) error {
	float64Lat, err := strconv.ParseFloat(request.Lat, 64)
	if err != nil {
		return fmt.Errorf("failed to convert latitude to float64: %w", err)
	}

	float64Lon, err := strconv.ParseFloat(request.Lon, 64)
	if err != nil {
		return fmt.Errorf("failed to convert longitude to float64: %w", err)
	}

	if err = s.repo.CreateLocation(userID, request.Name, float64Lat, float64Lon); err != nil {
		uniqErr := pq.As(err, pqerror.UniqueViolation)
		if uniqErr != nil {
			return ErrLocationAlreadyExists
		}

		return fmt.Errorf("failed to create location: %w", err)
	}

	return nil
}

func (s *Service) ListLocationsByUserID(userID int64) ([]*entity.Location, error) {
	locationsList, err := s.repo.ListLocationsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list locations: %w", err)
	}

	var result []*entity.Location

	for _, l := range locationsList {
		//	Wrapped into func so defer can immediately
		//	close resp.Body on the end of each iteration
		location, err := func() (*entity.Location, error) {
			stringLat := strconv.FormatFloat(l.Coordinates.Lat, 'f', -1, 64)
			stringLon := strconv.FormatFloat(l.Coordinates.Lon, 'f', -1, 64)

			resp, err := s.client.GetByCityNameAndCoordinates(l.Name, stringLat, stringLon)
			if err != nil {
				return nil, fmt.Errorf("failed to get locations by coordinates: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			}

			var location entity.Location
			if err = json.NewDecoder(resp.Body).Decode(&location); err != nil {
				return nil, fmt.Errorf("failed to decode location: %w", err)
			}

			//	Round values
			location.Main.Temp = math.Round(location.Main.Temp)
			location.Main.FeelsLike = math.Round(location.Main.FeelsLike)

			return &location, nil
		}()
		if err != nil {
			return nil, fmt.Errorf("failed to list locations: %w", err)
		}

		result = append(result, location)
	}

	return result, nil
}

func (s *Service) DeleteLocation(userID int64, name string) error {
	if err := s.repo.DeleteLocation(userID, name); err != nil {
		return fmt.Errorf("failed to delete location: %w", err)
	}

	return nil
}
