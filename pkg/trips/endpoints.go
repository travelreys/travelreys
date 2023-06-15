package trips

import (
	context "context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/media"
	"github.com/travelreys/travelreys/pkg/reqctx"
	"github.com/travelreys/travelreys/pkg/storage"
)

// Trips Endpoints

type CreateRequest struct {
	Name      string    `json:"name"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}
type CreateResponse struct {
	Trip Trip  `json:"trip"`
	Err  error `json:"error,omitempty"`
}

func (r CreateResponse) Error() error {
	return r.Err
}

func NewCreateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(CreateRequest)
		if !ok {
			return CreateResponse{
				Err: common.ErrorEndpointReqMismatch}, nil
		}

		ci, err := reqctx.ClientInfoFromCtx(ctx)
		if err != nil {
			return CreateResponse{Trip: Trip{}, Err: ErrRBAC}, nil
		}

		creator := NewMember(ci.UserID, MemberRoleCreator)
		trip, err := svc.Create(ctx, creator, req.Name, req.StartDate, req.EndDate)
		return CreateResponse{Trip: trip, Err: err}, nil
	}
}

type ReadRequest struct {
	ID          string `json:"id"`
	WithMembers bool   `json:"withMembers"`
}

type ReadResponse struct {
	Trip Trip  `json:"trip"`
	Err  error `json:"error,omitempty"`
}

func (r ReadResponse) Error() error {
	return r.Err
}

type ReadWithMembersResponse struct {
	Trip    Trip          `json:"trip"`
	Members auth.UsersMap `json:"members"`
	Err     error         `json:"error,omitempty"`
}

func (r ReadWithMembersResponse) Error() error {
	return r.Err
}

func NewReadEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ReadRequest)
		if !ok {
			return ReadResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		if req.WithMembers {
			trip, members, err := svc.ReadWithMembers(ctx, req.ID)
			return ReadWithMembersResponse{
				Trip: trip, Members: members, Err: err,
			}, nil
		}
		trip, err := svc.Read(ctx, req.ID)
		return ReadResponse{Trip: trip, Err: err}, nil
	}
}

type ReadMembersRequest struct {
	ID string `json:"id"`
}

type ReadMembersResponse struct {
	Members auth.UsersMap `json:"members"`
	Err     error         `json:"error,omitempty"`
}

func (r ReadMembersResponse) Error() error {
	return r.Err
}

func NewReadMembersEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ReadMembersRequest)
		if !ok {
			return ReadMembersResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		members, err := svc.ReadMembers(ctx, req.ID)
		return ReadMembersResponse{Members: members, Err: err}, nil
	}
}

type ReadOGPResponse struct {
	TripOGP TripOGP `json:"ogp"`
	Err     error   `json:"error,omitempty"`
}

func (r ReadOGPResponse) Error() error {
	return r.Err
}

func NewReadOGPEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ReadRequest)
		if !ok {
			return ReadOGPResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		ogp, err := svc.ReadOGP(ctx, req.ID)
		return ReadOGPResponse{TripOGP: ogp, Err: err}, nil
	}
}

type ListRequest struct {
	ListFilter
}
type ListResponse struct {
	Trips TripsList `json:"trips"`
	Err   error     `json:"error,omitempty"`
}

func (r ListResponse) Error() error {
	return r.Err
}

func NewListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListRequest)
		if !ok {
			return ListResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		trips, err := svc.List(ctx, req.ListFilter)
		return ListResponse{Trips: trips, Err: err}, nil
	}
}

type DeleteRequest struct {
	ID string `json:"id"`
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
		err := svc.Delete(ctx, req.ID)
		return DeleteResponse{Err: err}, nil
	}
}

type DeleteAttachmentRequest struct {
	ID  string         `json:"id"`
	Obj storage.Object `json:"object"`
}

type DeleteAttachmentResponse struct {
	Err error `json:"error,omitempty"`
}

func (r DeleteAttachmentResponse) Error() error {
	return r.Err
}

func NewDeleteAttachmentEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(DeleteAttachmentRequest)
		if !ok {
			return DeleteAttachmentResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		err := svc.DeleteAttachment(ctx, req.ID, req.Obj)
		return DeleteAttachmentResponse{Err: err}, nil
	}
}

type DownloadAttachmentPresignedURLRequest struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
}

type DownloadAttachmentPresignedURLResponse struct {
	PresignedURL string `json:"presignedURL"`
	Err          error  `json:"error,omitempty"`
}

func (r DownloadAttachmentPresignedURLResponse) Error() error {
	return r.Err
}

func NewDownloadAttachmentPresignedURLEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(DownloadAttachmentPresignedURLRequest)
		if !ok {
			return DownloadAttachmentPresignedURLResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		presignedURL, err := svc.DownloadAttachmentPresignedURL(ctx, req.ID, req.Path, req.Filename)
		return DownloadAttachmentPresignedURLResponse{
			PresignedURL: presignedURL,
			Err:          err,
		}, nil
	}
}

type UploadAttachmentPresignedURLRequest struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
}

type UploadAttachmentPresignedURLResponse struct {
	PresignedURL string `json:"presignedURL"`
	Err          error  `json:"error,omitempty"`
}

func (r UploadAttachmentPresignedURLResponse) Error() error {
	return r.Err
}

func NewUploadAttachmentPresignedURLEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(UploadAttachmentPresignedURLRequest)
		if !ok {
			return UploadAttachmentPresignedURLResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		presignedURL, err := svc.UploadAttachmentPresignedURL(ctx, req.ID, req.Filename)
		return UploadAttachmentPresignedURLResponse{
			PresignedURL: presignedURL,
			Err:          err,
		}, nil
	}
}

type GenerateMediaItemsRequest struct {
	ID     string                     `json:"id"`
	UserID string                     `json:"userID"`
	Params []media.NewMediaItemParams `json:"params"`
}

type GenerateMediaItemsResponse struct {
	Items media.MediaItemList         `json:"items"`
	URLs  media.MediaPresignedUrlList `json:"urls"`
	Err   error                       `json:"error,omitempty"`
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
		items, urls, err := svc.GenerateMediaItems(ctx, req.ID, req.UserID, req.Params)
		return GenerateMediaItemsResponse{Items: items, URLs: urls, Err: err}, nil
	}
}

type SaveMediaItemsRequest struct {
	ID    string              `json:"id"`
	Items media.MediaItemList `json:"items"`
}

type SaveMediaItemsResponse struct {
	Err error `json:"error,omitempty"`
}

func (r SaveMediaItemsResponse) Error() error {
	return r.Err
}

func NewSaveMediaItemsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(SaveMediaItemsRequest)
		if !ok {
			return SaveMediaItemsResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		err := svc.SaveMediaItems(ctx, req.ID, req.Items)
		return SaveMediaItemsResponse{Err: err}, nil
	}
}

type GenerateSignedURLsRequest struct {
	ID    string              `json:"id"`
	Items media.MediaItemList `json:"items"`
}

type GenerateSignedURLsResponse struct {
	URLs media.MediaPresignedUrlList `json:"urls"`
	Err  error                       `json:"error,omitempty"`
}

func (r GenerateSignedURLsResponse) Error() error {
	return r.Err
}

func NewGenerateSignedURLsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(GenerateSignedURLsRequest)
		if !ok {
			return GenerateSignedURLsResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		urls, err := svc.GenerateGetSignedURLs(ctx, req.ID, req.Items)
		return GenerateSignedURLsResponse{URLs: urls, Err: err}, nil
	}
}

type DeleteMediaItemsRequest struct {
	ID    string              `json:"id"`
	Items media.MediaItemList `json:"items"`
}

type DeleteMediaItemsResponse struct {
	Err error `json:"error,omitempty"`
}

func (r DeleteMediaItemsResponse) Error() error {
	return r.Err
}

func NewDeleteMediaItemsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(DeleteMediaItemsRequest)
		if !ok {
			return DeleteMediaItemsResponse{Err: common.ErrorEndpointReqMismatch}, nil
		}
		err := svc.DeleteMediaItems(ctx, req.ID, req.Items)
		return DeleteMediaItemsResponse{Err: err}, nil
	}
}
