package trips

import (
	context "context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/travelreys/travelreys/pkg/auth"
	"github.com/travelreys/travelreys/pkg/common"
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
		plan, err := svc.Create(ctx, creator, req.Name, req.StartDate, req.EndDate)
		return CreateResponse{Trip: plan, Err: err}, nil
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
			plan, members, err := svc.ReadWithMembers(ctx, req.ID)
			return ReadWithMembersResponse{
				Trip: plan, Members: members, Err: err,
			}, nil
		}
		plan, err := svc.Read(ctx, req.ID)
		return ReadResponse{Trip: plan, Err: err}, nil
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

type ListRequest struct {
	FF ListFilter
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
		plans, err := svc.List(ctx, req.FF)
		return ListResponse{Trips: plans, Err: err}, nil
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
