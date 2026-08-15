package view

import (
	"github.com/Nurlan270/weather-go/internal/rest/openweather/response"
	"github.com/Nurlan270/weather-go/internal/view"
)

type SearchViewData struct {
	view.BaseViewData
	Locations []*response.SearchLocation
}

func NewSearchViewData(title, login string, locations []*response.SearchLocation) *SearchViewData {
	return &SearchViewData{
		BaseViewData: *view.NewBaseViewData(title, login),
		Locations:    locations,
	}
}
