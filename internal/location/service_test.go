package location

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Nurlan270/weather-go/internal/cache"
	"github.com/Nurlan270/weather-go/internal/config"
	"github.com/Nurlan270/weather-go/internal/logger"
	"github.com/Nurlan270/weather-go/internal/rest/openweather/response"
	"github.com/Nurlan270/weather-go/internal/rest/request"
)

func TestLocationService(t *testing.T) {
	//	Arrange
	var (
		log       = logger.New(logger.EnvLocal)
		cacheConf = config.Cache{
			ExpiresIn: 5 * time.Minute,
		}
		c = cache.New(cacheConf, log)

		mockLocationRepo   = NewMockLocationRepository(t)
		mockLocationClient = NewMockLocationClient(t)

		locationSvc = NewService(mockLocationRepo, mockLocationClient, c)
	)

	//	Tests
	t.Run("it returns locations list by search", func(t *testing.T) {
		//	Arrange
		apiResp := getAPIResponseData()

		body, err := json.Marshal(apiResp)
		if err != nil {
			t.Fatalf("failed to serialize response body %s", err)
		}

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		}
		mockLocationClient.
			On("GetCitiesInfo", "London").
			Return(resp, nil)

		expected := []*response.SearchLocation{
			{
				Name:    "London",
				Country: "GB",
				Lat:     51.5073219,
				Lon:     -0.1276474,
			},
			{
				Name:    "London",
				Country: "CA",
				Lat:     42.9832406,
				Lon:     -81.243372,
			},
		}

		//	Act
		req := request.SearchLocation{
			Query: "London",
		}
		actual, err := locationSvc.SearchLocations(req)

		require.NoError(t, err, "Should not return an error")
		require.Len(t, actual, 2, "Should return response body of length 2")
		require.Equal(t, "GB", actual[0].Country, "Should return first city's country name correctly")
		require.Equal(t, expected, actual, "Should return response body as expected")
	})

	t.Run("it returns error if nothing was found by search query", func(t *testing.T) {
		//	Mocks
		mockLocationClient.
			On("GetCitiesInfo", "fictional country").
			Return(nil, ErrNoResults)

		//	Act
		req := request.SearchLocation{
			Query: "fictional country",
		}
		resp, err := locationSvc.SearchLocations(req)

		require.ErrorIs(t, err, ErrNoResults, "Should return ErrNoResults error")
		require.Nil(t, resp, "Should return nil body")
	})
}

func getAPIResponseData() []map[string]any {
	//	For simplicity, we'll return only 2 locations as response result
	return []map[string]any{
		{
			"name": "London",
			"local_names": map[string]string{
				"ha": "Landan",
				"tt": "Лондон",
				"lb": "London",
				"ce": "Лондон",
				"hu": "London",
				"it": "Londra",
				"tl": "Londres",
				"pl": "Londyn",
			},
			"lat":     51.5073219,
			"lon":     -0.1276474,
			"country": "GB",
			"state":   "England",
		},
		{
			"name": "London",
			"local_names": map[string]string{
				"el": "Λόντον",
				"fr": "London",
				"oj": "Baketigweyaang",
				"en": "London",
				"bn": "লন্ডন",
				"be": "Лондан",
				"ko": "런던",
				"he": "לונדון",
				"ru": "Лондон",
			},
			"lat":     42.9832406,
			"lon":     -81.243372,
			"country": "CA",
			"state":   "Ontario",
		},
	}
}
