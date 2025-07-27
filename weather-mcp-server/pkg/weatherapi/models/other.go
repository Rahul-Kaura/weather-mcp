package models

type Location struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
}
