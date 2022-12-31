package common

type Positioning struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`

	Continent string `json:"continent"`
	Country   string `json:"country"`
	City      string `json:"city"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`

	Labels Labels `json:"labels"`
}
