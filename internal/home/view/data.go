package view

import (
	"github.com/Nurlan270/weather-go/internal/entity"
	"github.com/Nurlan270/weather-go/internal/view"
)

type HomeViewData struct {
	view.BaseViewData
	Locations []*entity.Location
}

func NewHomeViewData(title, login string, locations []*entity.Location) *HomeViewData {
	return &HomeViewData{
		BaseViewData: *view.NewBaseViewData(title, login),
		Locations:    locations,
	}
}
