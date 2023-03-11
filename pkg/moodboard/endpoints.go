package moodboard

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/travelreys/travelreys/pkg/common"
)

type ReadAndCreateIfNotExistsRequest struct{}
type ReadAndCreateIfNotExistsResponse struct {
	Moodboard Moodboard `json:"moodboard"`
	Err       error     `json:"error,omitempty"`
}

func (r ReadAndCreateIfNotExistsResponse) Error() error {
	return r.Err
}

func NewReadAndCreateIfNotExistsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		mb, err := svc.ReadAndCreateIfNotExists(ctx, "")
		return ReadAndCreateIfNotExistsResponse{Moodboard: mb, Err: err}, nil
	}
}

type UpdateRequest struct {
	Title string `json:"title"`
}
type UpdateResponse struct {
	Err error `json:"error,omitempty"`
}

func (r UpdateResponse) Error() error {
	return r.Err
}

func NewUpdateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(UpdateRequest)
		if !ok {
			return UpdateResponse{
				Err: common.ErrorEndpointReqMismatch,
			}, nil
		}
		err := svc.Update(ctx, "", req.Title)
		return UpdateResponse{Err: err}, nil
	}
}

type AddPinRequest struct {
	Url string `json:"url"`
}
type AddPinResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

func (r AddPinResponse) Error() error {
	return r.Err
}

func NewAddPinEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(AddPinRequest)
		if !ok {
			return AddPinResponse{
				Err: common.ErrorEndpointReqMismatch,
			}, nil
		}
		id, err := svc.AddPin(ctx, "", req.Url)
		return AddPinResponse{ID: id, Err: err}, nil
	}
}

type UpdatePinRequest struct {
	PinID string `json:"pinID"`
	Notes string `json:"notes"`
}
type UpdatePinResponse struct {
	Err error `json:"error,omitempty"`
}

func (r UpdatePinResponse) Error() error {
	return r.Err
}

func NewUpdatePinEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(UpdatePinRequest)
		if !ok {
			return UpdatePinResponse{
				Err: common.ErrorEndpointReqMismatch,
			}, nil
		}
		err := svc.UpdatePin(ctx, "", req.PinID, req.Notes)
		return UpdatePinResponse{Err: err}, nil
	}
}

type DeletePinRequest struct {
	PinID string `json:"pinID"`
}
type DeletePinResponse struct {
	Err error `json:"error,omitempty"`
}

func (r DeletePinResponse) Error() error {
	return r.Err
}

func NewDeletePinEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(DeletePinRequest)
		if !ok {
			return DeletePinResponse{
				Err: common.ErrorEndpointReqMismatch,
			}, nil
		}
		err := svc.DeletePin(ctx, "", req.PinID)
		return DeletePinResponse{Err: err}, nil
	}
}
