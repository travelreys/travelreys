package media

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/travelreys/travelreys/pkg/common"
)

type GenerateMediaItemsRequest struct {
	UserID string               `json:"userID"`
	Params []NewMediaItemParams `json:"params"`
}

type GenerateMediaItemsResponse struct {
	Items MediaItemList `json:"items"`
	URLs  []string      `json:"urls"`
	Err   error         `json:"error,omitempty"`
}

func (r GenerateMediaItemsResponse) Error() error {
	return r.Err
}

func NewGenerateMediaItemsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(GenerateMediaItemsRequest)
		if !ok {
			return GenerateMediaItemsResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		items, urls, err := svc.GenerateMediaItems(ctx, req.UserID, req.Params)
		return GenerateMediaItemsResponse{Items: items, URLs: urls, Err: err}, nil
	}
}

type SaveForUserRequest struct {
	UserID string        `json:"userID"`
	Items  MediaItemList `json:"items"`
}

type SaveForUserResponse struct {
	Err error `json:"error,omitempty"`
}

func (r SaveForUserResponse) Error() error {
	return r.Err
}

func NewSaveForUserEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(SaveForUserRequest)
		if !ok {
			return SaveForUserResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		err := svc.SaveForUser(ctx, req.UserID, req.Items)
		return SaveForUserResponse{Err: err}, nil
	}
}

type ListRequest struct {
	ListMediaFilter
	WithURLs bool `json:"withURLs"`
}

type ListResponse struct {
	Items MediaItemList `json:"items"`
	Err   error         `json:"error,omitempty"`
}

func (r ListResponse) Error() error {
	return r.Err
}

type ListWithSignedURLsResponse struct {
	Items MediaItemList `json:"items"`
	URLs  []string      `json:"urls"`
	Err   error         `json:"error,omitempty"`
}

func (r ListWithSignedURLsResponse) Error() error {
	return r.Err
}

func NewListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListRequest)
		if !ok {
			return ListResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}

		if req.WithURLs {
			items, urls, err := svc.ListWithSignedURLs(ctx, req.ListMediaFilter)
			return ListWithSignedURLsResponse{items, urls, err}, nil
		}

		items, err := svc.List(ctx, req.ListMediaFilter)
		return ListResponse{items, err}, nil
	}
}

type DeleteRequest struct {
	DeleteMediaFilter
}

type DeleteResponse struct {
	Err error `json:"error,omitempty"`
}

func (r DeleteResponse) Error() error {
	return r.Err
}

func NewDeleteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(DeleteRequest)
		if !ok {
			return DeleteResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		err := svc.Delete(ctx, req.DeleteMediaFilter)
		return DeleteResponse{err}, nil
	}
}

type GenerateUploadPresignURLsRequest struct {
	UserID  string   `json:"userID"`
	FileIDs []string `json:"fileIDs"`
}

type GenerateUploadPresignURLsResponse struct {
	URLs []string `json:"urls"`
	Err  error    `json:"error,omitempty"`
}

func (r GenerateUploadPresignURLsResponse) Error() error {
	return r.Err
}

func NewGenerateUploadPresignURLsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(GenerateUploadPresignURLsRequest)
		if !ok {
			return GenerateUploadPresignURLsResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		urls, err := svc.GenerateUploadPresignURLs(ctx, req.UserID, req.FileIDs)
		return GenerateUploadPresignURLsResponse{urls, err}, nil
	}
}
