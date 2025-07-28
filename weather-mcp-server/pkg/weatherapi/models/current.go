package models

type Current struct {
	TempC      float64   `json:"temp_c"`
	TempF      float64   `json:"temp_f"`
	WindKph    float64   `json:"wind_kph"`
	WindMph    float64   `json:"wind_mph"`
	WindDir    string    `json:"wind_dir"`
	Humidity   int64     `json:"humidity"`
	FeelslikeC float64   `json:"feelslike_c"`
	FeelslikeF float64   `json:"feelslike_f"`
	Visibility float64   `json:"vis_km"`
	UV         float64   `json:"uv"`
	GustKph    float64   `json:"gust_kph"`
	PressureMb float64   `json:"pressure_mb"`
	Condition  Condition `json:"condition"`
}

type CurrentResponse struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}
