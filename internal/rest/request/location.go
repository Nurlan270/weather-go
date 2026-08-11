package request

type SearchLocation struct {
	Query string `json:"query" validate:"required"`
}

type AddLocation struct {
	Name string `json:"name" validate:"required,min=1,max=100,alphaspace"`
	Lat  string `json:"lat" validate:"required,latitude"`
	Lon  string `json:"lon" validate:"required,longitude"`
}
