package trips

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/reqctx"
	"github.com/travelreys/travelreys/pkg/storage"
)

const (
	URLPathVarID = "id"
)

func errToHttpCode(err error) int {
	notFoundErrors := []error{ErrTripNotFound}
	appErrors := []error{ErrUnexpectedStoreError}

	if common.ErrorContains(notFoundErrors, err) {
		return http.StatusNotFound
	}
	if common.ErrorContains(appErrors, err) {
		return http.StatusUnprocessableEntity
	}
	if errors.Is(err, ErrRBAC) {
		return http.StatusUnauthorized
	}
	if errors.Is(err, common.ErrValidation) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(common.Errorer); ok && e.Error() != nil {
		common.EncodeErrorFactory(errToHttpCode)(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Encoding", "gzip")

	gw := gzip.NewWriter(w)
	defer gw.Close()

	return json.NewEncoder(gw).Encode(response)
}

func MakeHandler(svc Service) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(reqctx.ContextWithClientInfo),
		kithttp.ServerErrorEncoder(common.EncodeErrorFactory(errToHttpCode)),
	}

	createHandler := kithttp.NewServer(
		NewCreateEndpoint(svc), decodeCreateRequest, encodeResponse, opts...,
	)
	listHandler := kithttp.NewServer(
		NewListEndpoint(svc), decodeListRequest, encodeResponse, opts...,
	)

	readHandler := kithttp.NewServer(
		NewReadEndpoint(svc), decodeReadRequest, encodeResponse, opts...,
	)

	readMembersHandler := kithttp.NewServer(
		NewReadMembersEndpoint(svc), decodeReadMembersRequest, encodeResponse, opts...,
	)
	readOGPHandler := kithttp.NewServer(
		NewReadOGPEndpoint(svc), decodeReadRequest, encodeResponse, opts...,
	)
	deleteHandler := kithttp.NewServer(
		NewDeleteEndpoint(svc), decodeDeleteRequest, encodeResponse, opts...,
	)
	deleteAttachmentHandler := kithttp.NewServer(
		NewDeleteAttachmentEndpoint(svc),
		decodeDeleteAttachmentRequest, encodeResponse, opts...,
	)
	uploadAttachmentPresignedURLHandler := kithttp.NewServer(
		NewUploadAttachmentPresignedURLEndpoint(svc),
		decodeUploadAttachmentPresignedURLRequest, encodeResponse, opts...,
	)
	downloadAttachmentPresignedURLHandler := kithttp.NewServer(
		NewDownloadAttachmentPresignedURLEndpoint(svc),
		decodeDownloadAttachmentPresignedURLRequest, encodeResponse, opts...,
	)

	generateMediaItemsHandler := kithttp.NewServer(
		NewGenerateMediaItemsEndpoint(svc),
		decodeGenerateMediaItemsRequest, encodeResponse, opts...,
	)

	saveMediaItemsHandler := kithttp.NewServer(
		NewSaveMediaItemsEndpoint(svc),
		decodeSaveMediaItemsRequest, encodeResponse, opts...,
	)

	deleteMediaItemsHandler := kithttp.NewServer(
		NewDeleteMediaItemsEndpoint(svc),
		decodeDeleteMediaItemsRequest, encodeResponse, opts...,
	)

	generateSignedURLsHandler := kithttp.NewServer(
		NewGenerateSignedURLsEndpoint(svc),
		decodeGenerateSignedURLsRequest, encodeResponse, opts...,
	)

	r.Handle("/api/v1/trips", createHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/trips", listHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}", readHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}/ogp", readOGPHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}/members", readMembersHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}", deleteHandler).Methods(http.MethodDelete)

	r.Handle("/api/v1/trips/{id}/storage", deleteAttachmentHandler).Methods(http.MethodDelete)
	r.Handle("/api/v1/trips/{id}/storage/download/pre-signed", downloadAttachmentPresignedURLHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}/storage/upload/pre-signed", uploadAttachmentPresignedURLHandler).Methods(http.MethodGet)

	r.Handle("/api/v1/trips/{id}/media/items", saveMediaItemsHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/trips/{id}/media/items", deleteMediaItemsHandler).Methods(http.MethodDelete)
	r.Handle("/api/v1/trips/{id}/media/items/generate", generateMediaItemsHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/trips/{id}/media/pre-signed", generateSignedURLsHandler).Methods(http.MethodPost)

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
	ff := ListFilter{UserID: &ci.UserID}
	req := ListRequest{ListFilter: ff}
	if r.URL.Query().Get("withMembers") == "true" {
		req.WithMembers = true
	}
	return req, nil
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	return DeleteRequest{ID}, nil
}

// Attachments

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

// Media

func decodeGenerateMediaItemsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}

	req := GenerateMediaItemsRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	req.ID = ID
	return req, nil
}

func decodeSaveMediaItemsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}

	req := SaveMediaItemsRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	req.ID = ID
	return req, nil
}

func decodeDeleteMediaItemsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}

	req := DeleteMediaItemsRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	req.ID = ID
	return req, nil
}

func decodeGenerateSignedURLsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}

	req := GenerateSignedURLsRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	req.ID = ID

	return req, nil
}
