package openweather

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/Nurlan270/weather-go/internal/config"
)

const (
	BaseURL           = "https://api.openweathermap.org"
	GeocodingEndpoint = BaseURL + "/geo/1.0/direct"
	WeatherEndpoint   = BaseURL + "/data/2.5/weather"

	//	Limit is a maximum number of locations that can be returned at once (max. 5)
	Limit = "5"
)

type Client struct {
	httpClient *http.Client
	apiCfg     config.OpenWeather
}

func NewClient(apiCfg config.OpenWeather) *Client {
	return &Client{
		httpClient: http.DefaultClient,
		apiCfg:     apiCfg,
	}
}

func (c *Client) GetCitiesInfo(cityName string) (*http.Response, error) {
	params := url.Values{}
	params.Set("q", cityName)
	params.Set("limit", Limit)

	endpoint := c.buildEndpoint(GeocodingEndpoint, params)

	resp, err := c.get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get %q: %w", endpoint, err)
	}

	return resp, nil
}

func (c *Client) GetCityWeather(lat, lon string) (*http.Response, error) {
	params := url.Values{}
	params.Set("lat", lat)
	params.Set("lon", lon)

	endpoint := c.buildEndpoint(WeatherEndpoint, params)

	resp, err := c.get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get %q: %w", endpoint, err)
	}

	return resp, nil
}

// get is low-level helper method.
func (c *Client) get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call request: %w", err)
	}

	return resp, nil
}

// buildEndpoint builds API's endpoint and pre-defines
// API Key and some other options on advance.
func (c *Client) buildEndpoint(endpoint string, params url.Values) string {
	const units = "metric"

	params.Set("appid", c.apiCfg.ApiKey)
	params.Set("units", units)

	return endpoint + "?" + params.Encode()
}
