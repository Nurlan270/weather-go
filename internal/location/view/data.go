package view

import (
	"github.com/Nurlan270/weather-go/internal/rest/openweather/response"
	"github.com/Nurlan270/weather-go/internal/view"
)

type SearchViewData struct {
	view.BaseViewData
	Location *response.Location
}

func NewSearchViewData(title, login string, location *response.Location) *SearchViewData {
	return &SearchViewData{
		BaseViewData: *view.NewBaseViewData(title, login),
		Location:     location,
	}
}
