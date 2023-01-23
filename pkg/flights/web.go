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
	"time"

	"github.com/tiinyplanet/tiinyplanet/pkg/common"
)

var (
	ErrInvalidSearchRequest = errors.New("invalid-flights-search-request")
)

// Flights Search

type FlightsSearchOptions map[string]string

type WebFlightsAPI interface {
	Search(
		ctx context.Context,
		origIATA,
		destIATA string,
		numAdults uint64,
		departDate time.Time,
		opts FlightsSearchOptions,
	) (ItinerariesList, error)
}

// Skyscanner

type SkyscannerResponse struct {
	Itineraries struct {
		Results []struct {
			ID   string `json:"id"`
			Legs []struct {
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
			} `json:"legs"`
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

func (res SkyscannerResponse) ToIinerariesList() ItinerariesList {
	list := ItinerariesList{}

	for _, result := range res.Itineraries.Results {
		metadata := result.Legs[0]
		pricing := result.PricingOptions[0]
		marketingAirline := metadata.Carriers.Marketing[0]
		itinerary := Itinerary{
			NumStops:           metadata.StopCount,
			Duration:           metadata.DurationInMinutes,
			Price:              pricing.Price.Amount,
			MarketingAirline:   Airline{Name: marketingAirline.Name},
			BookingURL:         pricing.URL,
			BookingDeeplinkURL: result.DeepLink,
			Legs:               LegsList{},
		}

		for idx, skyseg := range metadata.Segments {
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
			itinerary.Legs = append(itinerary.Legs, leg)

			if idx == 0 {
				itinerary.Departure = dep
				itinerary.MarketingAirline.Code = skyseg.MarketingCarrier.Code
			}
			if idx == len(metadata.Segments)-1 {
				itinerary.Arrival = arr
			}
		}
		list = append(list, itinerary)
	}
	return list
}

// Skyscanner API

type skyscanner struct {
	apiKey  string
	apiHost string
}

func NewSkyscannerAPI() WebFlightsAPI {
	apiKey := os.Getenv("TIINYPLANET_SKYSCANNER_APIKEY")
	apiHost := os.Getenv("TIINYPLANET_SKYSCANNER_APIHOST")
	return &skyscanner{apiKey, apiHost}
}

func (api skyscanner) fullUrlPath() string {
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

func (api skyscanner) Search(ctx context.Context, origIATA, destIATA string, numAdults uint64, departDate time.Time, opts FlightsSearchOptions) (ItinerariesList, error) {
	if origIATA == "" || destIATA == "" || numAdults <= 0 || departDate.Equal(time.Time{}) {
		return ItinerariesList{}, ErrInvalidSearchRequest
	}

	queryURL, _ := url.Parse(api.fullUrlPath())
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
	if val, ok := opts["currency"]; ok {
		queryParams.Set("currency", val)
	}
	if val, ok := opts["duration"]; ok {
		queryParams.Set("duration", val)
	}
	if val, ok := opts["stops"]; ok {
		queryParams.Set("stops", val)
	}
	queryURL.RawQuery = queryParams.Encode()

	body, err := api.Get(queryURL)
	if err != nil {
		return ItinerariesList{}, err
	}

	var res SkyscannerResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return ItinerariesList{}, err
	}

	return res.ToIinerariesList(), nil
}
