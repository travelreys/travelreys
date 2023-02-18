package flights

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tiinyplanet/tiinyplanet/pkg/common"
)

var (
	ErrInvalidSearchRequest = errors.New("invalid-flights-search-request")
)

const (
	DefaultCurrency = "USD"
)

type FlightsSearchOptions map[string]string

type WebFlightsAPI interface {
	Search(
		ctx context.Context,
		origIATA,
		destIATA string,
		numAdults uint64,
		departDate time.Time,
		opts FlightsSearchOptions,
	) (Itineraries, error)
}

// Skyscanner

type SkyscannerLeg struct {
	ID     string `json:"id"`
	Origin struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		DisplayCode string `json:"displayCode"`
	} `json:"origin"`
	Destination struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		DisplayCode string `json:"displayCode"`
	}
	DurationInMinutes uint64 `json:"durationInMinutes"`
	StopCount         uint64 `json:"stopCount"`
	IsSmallestStops   bool   `json:"isSmallestStops"`
	Departure         string `json:"departure"`
	Arrival           string `json:"arrival"`
	TimeDeltaInDays   uint32 `json:"timeDeltaInDays"`
	Carriers          struct {
		Marketing []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"marketing"`
		Operating []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"operating"`
		OperationType string `json:"opeartionType"`
	} `json:"carriers"`
	Segments []struct {
		ID     string `json:"id"`
		Origin struct {
			FlightPlaceID string `json:"flightPlaceId"`
			Name          string `json:"name"`
			Type          string `json:"type"`
			Parent        struct {
				FlightPlaceID string `json:"flightPlaceId"`
				Name          string `json:"string"`
				Type          string `json:"city"`
			} `json:"parent"`
		} `json:"origin"`
		Destination struct {
			FlightPlaceID string `json:"flightPlaceId"`
			Name          string `json:"name"`
			Type          string `json:"type"`
			Parent        struct {
				FlightPlaceID string `json:"flightPlaceId"`
				Name          string `json:"string"`
				Type          string `json:"city"`
			} `json:"parent"`
		} `json:"destination"`
		Departure         string `json:"departure"`
		Arrival           string `json:"arrival"`
		DurationInMinutes uint64 `json:"durationInMinutes"`
		FlightNumber      string `json:"flightNumber"`
		MarketingCarrier  struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Code string `json:"alternate_di"`
		} `json:"marketingCarrier"`
		OperatingCarrier struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Code string `json:"alternate_di"`
		} `json:"operatingCarrier"`
	} `json:"segments"`
}

type SkyscannerResponse struct {
	Itineraries struct {
		Results []struct {
			ID             string          `json:"id"`
			Legs           []SkyscannerLeg `json:"legs"`
			PricingOptions []struct {
				Agents []struct {
					ID        string `json:"id"`
					Name      string `json:"name"`
					IsCarrier bool   `json:"is_carrier"`
				} `json:"agents"`
				Price struct {
					Amount      float64 `json:"amount"`
					LastUpdated string  `json:"last_updated"`
				} `json:"price"`
				URL string `json:"url"`
			} `json:"pricing_options"`
			DeepLink string `json:"deeplink"`
		} `json:"results"`
	} `json:"itineraries"`
}

func calculateItineraryScore(price float64, durationInMins, stopCount uint64) float64 {
	return float64(price) + float64(durationInMins) + float64(60*stopCount)
}

func (res SkyscannerResponse) MakeFlightFromSkyLeg(skyleg SkyscannerLeg) Flight {
	flight := Flight{
		ID:       skyleg.ID,
		NumStops: skyleg.StopCount,
		Duration: skyleg.DurationInMinutes,
		Legs:     LegsList{},
	}
	for idx, skyseg := range skyleg.Segments {
		depDT, _ := time.Parse("2006-01-02T15:04:05", skyseg.Departure)
		arrDT, _ := time.Parse("2006-01-02T15:04:05", skyseg.Arrival)

		dep := Departure{
			Airport: Airport{
				Positioning: common.Positioning{
					Name: skyseg.Origin.Name,
				},
				Code: skyseg.Origin.FlightPlaceID,
			},
			Datetime: depDT,
		}
		arr := Arrival{
			Airport: Airport{
				Positioning: common.Positioning{
					Name: skyseg.Destination.Name,
				},
				Code: skyseg.Destination.FlightPlaceID,
			},
			Datetime: arrDT,
		}
		leg := Leg{
			FlightNo:  skyseg.FlightNumber,
			Departure: dep,
			Arrival:   arr,
			Duration:  skyseg.DurationInMinutes,
			OperatingAirline: Airline{
				Name: skyseg.OperatingCarrier.Name,
				Code: skyseg.OperatingCarrier.Code,
			},
		}
		flight.Legs = append(flight.Legs, leg)

		if idx == 0 {
			flight.Departure = dep
		}
		if idx == len(skyleg.Segments)-1 {
			flight.Arrival = arr
		}
	}
	return flight
}

func (res SkyscannerResponse) ToOnewayItineraries(currency string) Itineraries {
	itins := Itineraries{
		Type:    ITINERARY_ONE_WAY,
		Oneways: OnewayList{},
	}

	for _, result := range res.Itineraries.Results {
		departFlight := res.MakeFlightFromSkyLeg(result.Legs[0])
		pricing := result.PricingOptions[0]

		itin := Oneway{
			DepartFlight: departFlight,
			BookingMetadata: BookingMetadata{
				BookingURL: pricing.URL,
				Price: common.Price{
					Amount:   pricing.Price.Amount,
					Currency: currency,
				},
				BookingDeeplinkURL: result.DeepLink,
				Score: calculateItineraryScore(
					pricing.Price.Amount,
					departFlight.Duration,
					departFlight.NumStops),
			},
		}
		itins.Oneways = append(itins.Oneways, itin)
	}
	return itins
}

func (res SkyscannerResponse) ToRoundtripItineraries(currency string) Itineraries {
	itins := Itineraries{
		Type:       ITINERARY_ROUND_TRIP,
		Roundtrips: RoundTripMap{},
	}
	for _, result := range res.Itineraries.Results {
		departFlight := res.MakeFlightFromSkyLeg(result.Legs[0])
		returnFlight := res.MakeFlightFromSkyLeg(result.Legs[1])
		pricing := result.PricingOptions[0]

		bookingMetadata := BookingMetadata{
			BookingURL: pricing.URL,
			Price: common.Price{
				Amount:   pricing.Price.Amount,
				Currency: currency,
			},
			BookingDeeplinkURL: result.DeepLink,
			Score: calculateItineraryScore(
				pricing.Price.Amount, departFlight.Duration+
					returnFlight.Duration, departFlight.NumStops+
					returnFlight.NumStops),
		}
		if _, ok := itins.Roundtrips[departFlight.ID]; !ok {
			itins.Roundtrips[departFlight.ID] = &RoundTrip{
				DepartFlight:        departFlight,
				ReturnFlights:       FlightsList{returnFlight},
				BookingMetadataList: BookingMetadataList{bookingMetadata},
			}
		} else {
			itins.Roundtrips[departFlight.ID].ReturnFlights = append(
				itins.Roundtrips[departFlight.ID].ReturnFlights, returnFlight)
			itins.Roundtrips[departFlight.ID].BookingMetadataList = append(
				itins.Roundtrips[departFlight.ID].BookingMetadataList, bookingMetadata)
		}
	}
	return itins
}

// Skyscanner API

type skyscanner struct {
	apiKey  string
	apiHost string
}

func NewSkyscannerAPI(key, host string) WebFlightsAPI {
	return &skyscanner{key, host}
}

func (api skyscanner) urlpath() string {
	return fmt.Sprintf("https://%s/search-extended", api.apiHost)
}

func (api skyscanner) Get(getURL *url.URL) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, getURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-RapidAPI-Key", api.apiKey)
	req.Header.Add("X-RapidAPI-Host", api.apiHost)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	return body, err
}

func (api skyscanner) isRoundTrip(opts FlightsSearchOptions) bool {
	_, ok := opts["returnDate"]
	return ok
}

func (api skyscanner) Search(ctx context.Context, origIATA, destIATA string, numAdults uint64, departDate time.Time, opts FlightsSearchOptions) (Itineraries, error) {
	if origIATA == "" || destIATA == "" || numAdults <= 0 || departDate.Equal(time.Time{}) {
		return Itineraries{}, ErrInvalidSearchRequest
	}

	queryURL, _ := url.Parse(api.urlpath())
	queryParams := queryURL.Query()

	queryParams.Set("adults", fmt.Sprintf("%d", numAdults))
	queryParams.Set("origin", origIATA)
	queryParams.Set("destination", destIATA)
	queryParams.Set("departureDate", departDate.Format("2006-01-02"))

	if val, ok := opts["returnDate"]; ok {
		queryParams.Set("returnDate", val)
	}
	if val, ok := opts["cabinClass"]; ok {
		queryParams.Set("cabinClass", val)
	}

	currency := DefaultCurrency
	if val, ok := opts["currency"]; ok {
		currency = val
	}
	queryParams.Set("currency", currency)
	if val, ok := opts["duration"]; ok {
		queryParams.Set("duration", val)
	}
	if val, ok := opts["stops"]; ok {
		queryParams.Set("stops", val)
	}
	queryURL.RawQuery = queryParams.Encode()

	body, err := api.Get(queryURL)
	if err != nil {
		return Itineraries{}, err
	}

	var res SkyscannerResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return Itineraries{}, err
	}

	isRoundTrip := api.isRoundTrip(opts)
	if isRoundTrip {
		return res.ToRoundtripItineraries(currency), nil
	} else {
		return res.ToOnewayItineraries(currency), nil
	}
}
