package response

type Location struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	Coordinates Coordinates `json:"coord"`
	Sys         Sys         `json:"sys"`
	Main        Main        `json:"main"`
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Sys struct {
	Country string `json:"country"`
}

type Main struct {
	Temp     float64 `json:"temp"`
	Humidity int     `json:"humidity"`
}
