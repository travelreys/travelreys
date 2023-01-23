package common

type Positioning struct {
	Name    string `json:"name"`
	Address string `json:"address"`

	Continent string `json:"continent"`
	Country   string `json:"country"`
	State     string `json:"state"`
	City      string `json:"city"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`

	Labels Labels `json:"labels"`
}
