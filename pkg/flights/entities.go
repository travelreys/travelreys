package flights

import (
	"time"

	"github.com/tiinyplanet/tiinyplanet/pkg/common"
)

type Airline struct {
	Name        string `json:"name"`
	Code        string `json:"code"`       // IATA Code
	WebsiteURL  string `json:"websiteURL"` // May not be valid URL
	PhoneNumber string `json:"phoneNumber"`
}

type AirlinesList []Airline

type Airport struct {
	common.Positioning
	Code string `json:"code"`
}

type AirportsList []Airport

const (
	ITINERARY_ROUND_TRIP = "ROUND_TRIP"
	ITINERARY_ONE_WAY    = "ONE_WAY"
)

const (
	CABIN_CLASS_ECONOMY         = "economy"
	CABIN_CLASS_PREMIUM_ECONOMY = "premiumeconomy"
	CABIN_CLASS_BUSINESS        = "business"
	CABIN_CLASS_FIRST           = "first"
)

type Itinerary struct {
	Departure Departure `json:"departure"` // Initial departure
	Arrival   Arrival   `json:"arrival"`   // Final arrival
	NumStops  uint64    `json:"numStops"`  // Total number of stops
	Duration  uint64    `json:"duration"`  // Total duration in mins
	Price     float64   `json:"price"`
	Score     float64   `json:"score"` // ranking score = price + duration + 60 * numStops

	MarketingAirline   Airline `json:"marketingAirline"`   // Marketing Airline (sells the ticket)
	BookingURL         string  `json:"bookingURL"`         // URL to book the ticket
	BookingDeeplinkURL string  `json:"bookingDeeplinkURL"` // URL to see other options!

	Legs LegsList `json:"legs"`
}

type ItinerariesList []Itinerary

// A segment is a flight operated by a single flight number, but may have an intermediate stop
// Example - UA 234 from BOS-ORD-SFO is a segment.
type Segment struct{}

// A leg is always a single non-stop flight. Example, UA123 from BOS-EWR is a leg.
type Leg struct {
	FlightNo         string    `json:"flightNo"`
	Departure        Departure `json:"departure"`
	Arrival          Arrival   `json:"arrival"`
	Duration         uint64    `json:"duration"` // duration in mins
	OperatingAirline Airline   `json:"operatingAirline"`
}

type LegsList []Leg

type Departure struct {
	Airport  Airport   `json:"airport"`
	Datetime time.Time `json:"datetime"` // UTC
}

type Arrival struct {
	Airport  Airport   `json:"airport"`
	Datetime time.Time `json:"datetime"` // UTC
}
