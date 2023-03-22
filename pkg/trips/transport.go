package trips

import (
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/reqctx"
	"github.com/travelreys/travelreys/pkg/storage"
)

const (
	URLPathVarID  = "id"
	maxUploadSize = 25 * 1024 * 1024 // 25MB
)

func errToHttpCode(err error) int {
	notFoundErrors := []error{ErrPlanNotFound}
	appErrors := []error{ErrUnexpectedStoreError}
	authErrors := []error{ErrRBAC}

	if common.ErrorContains(notFoundErrors, err) {
		return http.StatusNotFound
	}
	if common.ErrorContains(appErrors, err) {
		return http.StatusUnprocessableEntity
	}
	if common.ErrorContains(authErrors, err) {
		return http.StatusUnauthorized
	}
	return http.StatusInternalServerError
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(common.Errorer); ok && e.Error() != nil {
		common.EncodeErrorFactory(errToHttpCode)(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func MakeHandler(svc Service) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(reqctx.ContextWithClientInfo),
		kithttp.ServerErrorEncoder(common.EncodeErrorFactory(errToHttpCode)),
	}

	createHandler := kithttp.NewServer(NewCreateEndpoint(svc), decodeCreateRequest, encodeResponse, opts...)
	listHandler := kithttp.NewServer(NewListEndpoint(svc), decodeListRequest, encodeResponse, opts...)
	readHandler := kithttp.NewServer(NewReadEndpoint(svc), decodeReadRequest, encodeResponse, opts...)
	readMembersHandler := kithttp.NewServer(NewReadMembersEndpoint(svc), decodeReadMembersRequest, encodeResponse, opts...)
	deleteHandler := kithttp.NewServer(NewDeleteEndpoint(svc), decodeDeleteRequest, encodeResponse, opts...)

	deleteAttachmentHandler := kithttp.NewServer(NewDeleteAttachmentEndpoint(svc), decodeDeleteAttachmentRequest, encodeResponse, opts...)
	downloadAttachmentPresignedURLHandler := kithttp.NewServer(NewDownloadAttachmentPresignedURLEndpoint(svc), decodeDownloadAttachmentPresignedURLRequest, encodeResponse, opts...)
	uploadAttachmentPresignedURLHandler := kithttp.NewServer(NewUploadAttachmentPresignedURLEndpoint(svc), decodeUploadAttachmentPresignedURLRequest, encodeResponse, opts...)
	uploadMediaPresignedURLHandler := kithttp.NewServer(NewUploadMediaPresignedURLEndpoint(svc), decodeUploadMediaPresignedURLRequest, encodeResponse, opts...)

	r.Handle("/api/v1/trips", createHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/trips", listHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}", readHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}/members", readMembersHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}", deleteHandler).Methods(http.MethodDelete)

	r.Handle("/api/v1/trips/{id}/storage", deleteAttachmentHandler).Methods(http.MethodDelete)
	r.Handle("/api/v1/trips/{id}/storage/download/pre-signed", downloadAttachmentPresignedURLHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}/storage/upload/pre-signed", uploadAttachmentPresignedURLHandler).Methods(http.MethodGet)

	r.Handle("/api/v1/trips/{id}/media/upload/pre-signed", uploadMediaPresignedURLHandler).Methods(http.MethodGet)

	// downloadHandler := kithttp.NewServer(NewDownloadEndpoint(svc), decodeDownloadRequest, encodeDownloadResponse, opts...)
	// r.Handle("/api/v1/trips/{id}/storage", downloadHandler).Methods(http.MethodGet)

	return r
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := CreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	return req, nil
}

func decodeReadRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	req := ReadRequest{ID: ID}
	if r.URL.Query().Get("withMembers") == "true" {
		req.WithMembers = true
	}
	return req, nil
}

func decodeReadMembersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	return ReadMembersRequest{ID: ID}, nil
}

func decodeListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	ci, err := reqctx.ClientInfoFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if ci.UserID == "" {
		return nil, ErrRBAC
	}
	ff := ListFilter{UserID: common.StringPtr(ci.UserID)}
	return ListRequest{ff}, nil

}
func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	return DeleteRequest{ID}, nil
}

func decodeDeleteAttachmentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	path := r.URL.Query().Get("path")
	return DeleteAttachmentRequest{
		ID: ID,
		Obj: storage.Object{
			Path: path,
			Name: filepath.Base(path),
		},
	}, nil
}

func decodeDownloadAttachmentPresignedURLRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	path := r.URL.Query().Get("path")
	filename := r.URL.Query().Get("filename")
	return DownloadAttachmentPresignedURLRequest{
		ID:       ID,
		Path:     path,
		Filename: filename,
	}, nil
}

func decodeUploadAttachmentPresignedURLRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	filename := r.URL.Query().Get("filename")
	return UploadAttachmentPresignedURLRequest{
		ID:       ID,
		Filename: filename,
	}, nil
}

func decodeUploadMediaPresignedURLRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	filename := r.URL.Query().Get("filename")
	return UploadMediaPresignedURLRequest{
		ID:       ID,
		Filename: filename,
	}, nil
}
