package maps

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
)

type PlacesAutocompleteRequest struct {
	Query        string `json:"query"`
	Types        string `json:"types"`
	Sessiontoken string `json:"sessiontoken"`
}
type PlacesAutocompleteResponse struct {
	Predictions AutocompletePredictionList `json:"predictions"`
	Err         error                      `json:"error,omitempty"`
}

func (r PlacesAutocompleteResponse) Error() error {
	return r.Err
}

func NewPlacesAutocompleteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(PlacesAutocompleteRequest)
		if !ok {
			return PlacesAutocompleteResponse{
				Err: common.ErrorEndpointReqMismatch,
			}, nil
		}
		preds, err := svc.PlacesAutocomplete(ctx, req.Query, req.Types, req.Sessiontoken)
		return PlacesAutocompleteResponse{Predictions: preds, Err: err}, nil
	}
}

type PlaceDetailsRequest struct {
	PlaceID      string   `json:"placeID"`
	Fields       []string `json:"fields"`
	Sessiontoken string   `json:"sessiontoken"`
}

type PlaceDetailsResponse struct {
	Place Place `json:"place"`
	Err   error `json:"error,omitempty"`
}

func (r PlaceDetailsResponse) Error() error {
	return r.Err
}

func NewPlaceDetailsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(PlaceDetailsRequest)
		if !ok {
			return PlaceDetailsResponse{
				Err: common.ErrorEndpointReqMismatch,
			}, nil
		}
		place, err := svc.PlaceDetails(ctx, req.PlaceID, req.Fields, req.Sessiontoken)
		return PlaceDetailsResponse{Place: place, Err: err}, nil
	}
}

type DirectionsRequest struct {
	OriginPlaceID string `json:"originPlaceID"`
	DestPlaceID   string `json:"destPlaceID"`
	Mode          string `json:"mode"`
}
type DirectionsResponse struct {
	RouteList RouteList `json:"routeList"`
	Err       error     `json:"error,omitempty"`
}

func (r DirectionsResponse) Error() error {
	return r.Err
}

func NewDirectionsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(DirectionsRequest)
		if !ok {
			return DirectionsResponse{
				Err: common.ErrorEndpointReqMismatch,
			}, nil
		}
		routeList, err := svc.Directions(ctx, req.OriginPlaceID, req.DestPlaceID, req.Mode)
		return DirectionsResponse{RouteList: routeList, Err: err}, nil
	}
}

type OptimizeRouteRequest struct {
	OriginPlaceID    string   `json:"originPlaceID"`
	DestPlaceID      string   `json:"destPlaceID"`
	WaypointsPlaceID []string `json:"waypointsPlaceID"`
}

type OptimizeRouteResponse struct {
	RouteList RouteList `json:"routeList"`
	Err       error     `json:"error,omitempty"`
}

func (r OptimizeRouteResponse) Error() error {
	return r.Err
}

func NewOptimizeRouteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(OptimizeRouteRequest)
		if !ok {
			return OptimizeRouteResponse{
				Err: common.ErrorEndpointReqMismatch,
			}, nil
		}
		routeList, err := svc.OptimizeRoute(ctx, req.OriginPlaceID, req.DestPlaceID, req.WaypointsPlaceID)
		return OptimizeRouteResponse{RouteList: routeList, Err: err}, nil
	}
}
