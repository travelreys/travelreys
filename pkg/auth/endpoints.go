package auth

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/travelreys/travelreys/pkg/common"
)

type MagicLinkRequest struct {
	Email string `json:"email"`
}

type MagicLinkResponse struct {
	Err error `json:"error,omitempty"`
}

func (r MagicLinkResponse) Error() error {
	return r.Err
}

func NewMagicLinkEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(MagicLinkRequest)
		if !ok {
			return MagicLinkResponse{Err: common.ErrEndpointReqMismatch}, nil
		}
		err := svc.MagicLink(ctx, req.Email)
		return MagicLinkResponse{Err: err}, nil
	}
}

type LoginRequest struct {
	Code      string `json:"code"`
	Signature string `json:"signature"`
	Provider  string `json:"provider"`
}

type LoginResponse struct {
	User   User         `json:"user"`
	Cookie *http.Cookie `json:"-"`
	Err    error        `json:"error,omitempty"`
}

func (r LoginResponse) Error() error {
	return r.Err
}

func NewLoginEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(LoginRequest)
		if !ok {
			return LoginResponse{Err: common.ErrEndpointReqMismatch}, nil
		}
		usr, cookie, err := svc.Login(ctx, req.Code, req.Signature, req.Provider)
		return LoginResponse{User: usr, Cookie: cookie, Err: err}, nil
	}
}

type ReadRequest struct {
	ID string `json:"id"`
}

type ReadResponse struct {
	User User  `json:"user"`
	Err  error `json:"error,omitempty"`
}

func (r ReadResponse) Error() error {
	return r.Err
}

func NewReadEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ReadRequest)
		if !ok {
			return ReadResponse{
				Err: common.ErrEndpointReqMismatch,
			}, nil
		}
		usr, err := svc.Read(ctx, req.ID)
		return ReadResponse{usr, err}, nil
	}
}

type UpdateRequest struct {
	ID string       `json:"id"`
	FF UpdateFilter `json:"ff"`
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
				Err: common.ErrEndpointReqMismatch,
			}, nil
		}
		err := svc.Update(ctx, req.ID, req.FF)
		return UpdateResponse{err}, nil
	}
}

type ListRequest struct {
	FF ListFilter
}
type ListResponse struct {
	Users UsersList `json:"users"`
	Err   error     `json:"error,omitempty"`
}

func (r ListResponse) Error() error {
	return r.Err
}

func NewListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(ListRequest)
		if !ok {
			return ListResponse{Err: common.ErrEndpointReqMismatch}, nil
		}
		users, err := svc.List(ctx, req.FF)
		return ListResponse{Users: users, Err: err}, nil
	}
}

type LogoutRequest struct{}

type LogoutResponse struct {
	Err error `json:"error,omitempty"`
}

func (r LogoutResponse) Error() error {
	return r.Err
}

func NewLogoutEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		return LogoutResponse{}, nil
	}
}

type DeleteRequest struct {
	ID string
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
			return DeleteResponse{
				Err: common.ErrEndpointReqMismatch,
			}, nil
		}
		err := svc.Delete(ctx, req.ID)
		return DeleteResponse{err}, nil
	}
}

type UploadAvatarPresignedURLRequest struct {
	ID       string
	MIMEType string
}

type UploadAvatarPresignedURLResponse struct {
	PresignedURL string `json:"presignedURL"`
	SuffixToken  string `json:"suffixToken"`
	Err          error  `json:"error,omitempty"`
}

func (r UploadAvatarPresignedURLResponse) Error() error {
	return r.Err
}

func NewUploadAvatarPresignedURLEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, epReq interface{}) (interface{}, error) {
		req, ok := epReq.(UploadAvatarPresignedURLRequest)
		if !ok {
			return UploadAvatarPresignedURLResponse{Err: common.ErrEndpointReqMismatch}, nil
		}
		suffixToken, presignedURL, err := svc.UploadAvatarPresignedURL(ctx, req.ID, req.MIMEType)
		return UploadAvatarPresignedURLResponse{
			PresignedURL: presignedURL,
			SuffixToken:  suffixToken,
			Err:          err,
		}, nil
	}
}
