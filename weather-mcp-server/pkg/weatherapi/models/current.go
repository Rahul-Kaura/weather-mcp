package models

type Current struct {
	TempC     float64   `json:"temp_c"`
	WindKph   float64   `json:"wind_kph"`
	Humidity  int64     `json:"humidity"`
	Condition Condition `json:"condition"`
}

type CurrentResponse struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}
