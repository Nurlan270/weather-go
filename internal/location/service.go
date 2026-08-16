package location

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/allegro/bigcache"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"

	"github.com/Nurlan270/weather-go/internal/cache"
	"github.com/Nurlan270/weather-go/internal/entity"
	"github.com/Nurlan270/weather-go/internal/rest/openweather/response"
	"github.com/Nurlan270/weather-go/internal/rest/request"
)

type LocationRepository interface {
	ListLocationIdsByUserID(userID int64) ([]int64, error)
	GetLocationByID(id int64) (entity.Location, error)
	CreateLocation(userID int64, name string, lat, lon float64) error
	DeleteLocation(id, userID int64) error
}

type LocationClient interface {
	GetCitiesInfo(cityName string) (*http.Response, error)
	GetCityWeather(lat, lon string) (*http.Response, error)
}

type Service struct {
	repo   LocationRepository
	client LocationClient
	cache  *cache.Cache
}

func NewService(repo LocationRepository, client LocationClient, cache *cache.Cache) *Service {
	return &Service{
		repo:   repo,
		client: client,
		cache:  cache,
	}
}

func (s *Service) SearchLocations(request request.SearchLocation) ([]*response.SearchLocation, error) {
	resp, err := s.client.GetCitiesInfo(request.Query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unepected status code %d", resp.StatusCode)
	}

	var locations []*response.SearchLocation
	if err = json.NewDecoder(resp.Body).Decode(&locations); err != nil {
		return nil, fmt.Errorf("failed to decode locations: %w", err)
	}

	//	Nothing found
	if len(locations) == 0 {
		return nil, ErrNoResults
	}

	return locations, err
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

func (s *Service) ListLocationsByUser(user *entity.User) ([]*entity.Location, error) {
	// 1. Load IDs from DB
	locationIDs, err := s.repo.ListLocationIdsByUserID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get locations from db: %w", err)
	}

	result := make([]*entity.Location, 0, len(locationIDs))

	for _, locationID := range locationIDs {
		locationKey := cache.BuildKey(cache.KeyUserLocations, user.ID, locationID)

		var location entity.Location
		if err = s.cache.GetInto(locationKey, &location); err == nil {
			//	2. Cache hit: append location to result, otherwise skip
			result = append(result, &location)
			continue
		}

		if !errors.Is(err, bigcache.ErrEntryNotFound) {
			return nil, fmt.Errorf("failed to get location from cache: %w", err)
		}

		//	3. Cache miss: get location from DB
		location, err = s.repo.GetLocationByID(locationID)
		if err != nil {
			return nil, fmt.Errorf("failed to get location by id %d: %w", locationID, err)
		}

		//	4. Cache miss: call client and put data into cache
		location, err = s.fetchAndCacheLocation(user.ID, location)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch and cache location by id %d: %w", locationID, err)
		}

		result = append(result, &location)
	}

	return result, nil
}

func (s *Service) DeleteLocation(id int64, user *entity.User) error {
	if err := s.repo.DeleteLocation(id, user.ID); err != nil {
		return fmt.Errorf("failed to delete location from db: %w", err)
	}

	key := cache.BuildKey(cache.KeyUserLocations, user.ID, id)
	if err := s.cache.Del(key); err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
		return fmt.Errorf("failed to delete location from cache: %w", err)
	}

	return nil
}

func (s *Service) fetchAndCacheLocation(userID int64, location entity.Location) (entity.Location, error) {
	lat := strconv.FormatFloat(location.Coordinates.Lat, 'f', -1, 64)
	lon := strconv.FormatFloat(location.Coordinates.Lon, 'f', -1, 64)

	//	Call API to get fresh weather info
	resp, err := s.client.GetCityWeather(lat, lon)
	if err != nil {
		return entity.Location{}, fmt.Errorf("failed to fetch location: %w", err)
	}
	defer resp.Body.Close()

	var result entity.Location

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return entity.Location{}, fmt.Errorf("failed to deserialize location: %w", err)
	}

	//	Re-assign ID from API to this app's ID
	//	It's then used to delete location
	result.ID = location.ID

	//	Re-assign name to name from DB
	result.Name = location.Name

	//	Round temperature
	result.Main.Temp = math.Round(result.Main.Temp)
	result.Main.FeelsLike = math.Round(result.Main.FeelsLike)

	locationKey := cache.BuildKey(cache.KeyUserLocations, userID, location.ID)

	//	Cache location
	if err = s.cache.Set(locationKey, result); err != nil {
		return entity.Location{}, fmt.Errorf("failed to cache location: %w", err)
	}

	return result, nil
}
