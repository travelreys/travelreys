package flights

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"go.uber.org/zap"
)

const (
	DefaultCurrency      = "USD"
	envSkyscannerApiKey  = "TIINYPLANET_SKYSCANNER_APIKEY"
	envSkyscannerApiHost = "TIINYPLANET_SKYSCANNER_APIHOST"

	fieldAdults        = "adults"
	fieldOrigin        = "origin"
	fieldDestination   = "destination"
	fieldDepartureDate = "departureDate"
	fieldCurrency      = "currency"
	fieldCabinClass    = "cabin"
	fieldReturnData    = "returnDate"
	fieldDuration      = "duration"
	fieldStops         = "stops"

	webLoggerName = "flights.webapi"
)

var (
	ErrInvalidSearchRequest = errors.New("web.invalidrequests")
)

type SearchOptions struct {
	Adults        int64
	Origin        string
	Destination   string
	DepartureDate time.Time
	ReturnDate    time.Time
	Currency      string
	CabinClass    string
	Duration      string
	Stops         string
}

func SearchOptionsFromURLValues(q url.Values) (SearchOptions, error) {
	opts := SearchOptions{}
	if q.Has(fieldCabinClass) {
		opts.CabinClass = q.Get(fieldCabinClass)
	}
	if q.Has(fieldCurrency) {
		opts.Currency = q.Get(fieldCurrency)
	}
	if q.Has(fieldDuration) {
		opts.Duration = q.Get(fieldDuration)
	}
	if q.Has(fieldStops) {
		opts.Stops = q.Get(fieldStops)
	}
	if q.Has(fieldReturnData) {
		dt, err := time.Parse("2006-01-02", q.Get(fieldReturnData))
		if err != nil {
			opts.ReturnDate = time.Time{}
		}
		opts.ReturnDate = dt
	}

	dt, err := time.Parse("2006-01-02", q.Get("departDate"))
	if err != nil {
		return opts, err
	}

	adults, err := strconv.Atoi(q.Get(fieldAdults))
	if err != nil {
		return opts, err
	}
	opts.Origin = q.Get(fieldOrigin)
	opts.Destination = q.Get(fieldDestination)
	opts.DepartureDate = dt
	opts.Adults = int64(adults)
	return opts, nil
}

func (opts SearchOptions) Validate() error {
	if (opts.Origin == "" ||
		opts.Destination == "" ||
		opts.Adults <= 0 ||
		opts.DepartureDate.Equal(time.Time{})) {
		return ErrInvalidSearchRequest
	}
	return nil
}

func (opts SearchOptions) IsRoundTrip() bool {
	return !opts.ReturnDate.Equal(time.Time{})
}

func (opts SearchOptions) ToRawURLQuery() string {
	values := url.Values{}

	values.Set(fieldOrigin, opts.Origin)
	values.Set(fieldDestination, opts.Destination)
	values.Set(fieldAdults, fmt.Sprintf("%d", opts.Adults))
	values.Set(fieldDepartureDate, opts.DepartureDate.Format("2006-01-02"))

	if !opts.ReturnDate.Equal(time.Time{}) {
		values.Set(fieldReturnData, opts.ReturnDate.Format("2006-01-02"))
	}
	if opts.CabinClass != "" {
		values.Set(fieldCabinClass, opts.CabinClass)
	}
	if opts.Duration != "" {
		values.Set(fieldDuration, opts.Duration)
	}
	if opts.Stops != "" {
		values.Set(fieldStops, opts.Stops)
	}
	if opts.Currency != "" {
		values.Set(fieldCurrency, opts.Currency)
	}

	return values.Encode()
}

type WebAPI interface {
	Search(ctx context.Context, opts SearchOptions) (Itineraries, error)
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

	logger *zap.Logger
}

func GetApiHost() string {
	return os.Getenv(envSkyscannerApiHost)
}

func GetApiKey() string {
	return os.Getenv(envSkyscannerApiKey)
}

func NewDefaultWebAPI(logger *zap.Logger) WebAPI {
	return NewWebAPI(GetApiHost(), GetApiKey(), logger)
}

func NewWebAPI(key, host string, logger *zap.Logger) WebAPI {
	return &skyscanner{key, host, logger.Named(webLoggerName)}
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
	return io.ReadAll(resp.Body)
}

func (api skyscanner) Search(ctx context.Context, opts SearchOptions) (Itineraries, error) {
	queryURL, _ := url.Parse(fmt.Sprintf("https://%s/search-extended", api.apiHost))
	queryURL.RawQuery = opts.ToRawURLQuery()

	currency := opts.Currency
	if currency == "" {
		currency = DefaultCurrency
	}

	body, err := api.Get(queryURL)
	if err != nil {
		api.logger.Error("Search", zap.Error(err))
		return Itineraries{}, err
	}

	var res SkyscannerResponse
	if err = json.Unmarshal(body, &res); err != nil {
		api.logger.Error("Search", zap.Error(err))
		return Itineraries{}, err
	}

	if opts.IsRoundTrip() {
		return res.ToRoundtripItineraries(currency), nil
	}
	return res.ToOnewayItineraries(currency), nil
}
