package request

type SearchLocation struct {
	Query string `json:"query" validate:"required,max=70"`
}

type AddLocation struct {
	Name string `json:"name" validate:"required,max=70"`
	Lat  string `json:"lat"  validate:"required,latitude"`
	Lon  string `json:"lon"  validate:"required,longitude"`
}
