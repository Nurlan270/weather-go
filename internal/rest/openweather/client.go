package openweather

import (
	"fmt"
	"net/http"

	"github.com/Nurlan270/weather-go/internal/config"
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

func (c *Client) GetByCityName(query string) (*http.Response, error) {
	endpoint := c.buildEndpoint(getByCityName, query)

	resp, err := c.get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get %q: %w", endpoint, err)
	}

	return resp, nil
}

func (c *Client) GetByCityNameAndCoordinates(name string, lat, lon float64) (*http.Response, error) {
	endpoint := c.buildEndpoint(getByCityNameAndCoordinates, name, lat, lon)

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
// API Key and some other options on advance
func (c *Client) buildEndpoint(endpoint string, a ...any) string {
	const units = "metric"

	args := append([]any{c.apiCfg.ApiKey, units}, a...)

	return fmt.Sprintf(endpoint, args...)
}
