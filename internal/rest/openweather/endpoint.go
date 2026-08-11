package openweather

const BASE = "https://api.openweathermap.org/data/2.5/weather?appid=%s&units=%s"

var (
	getByCityName               = BASE + "&q=%s"
	getByCityNameAndCoordinates = BASE + "&q=%s&lat=%f&lon=%f"
)
