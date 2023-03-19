package trips

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	uploadHandler := kithttp.NewServer(NewUploadEndpoint(svc), decodeUploadRequest, encodeUploadResponse, opts...)
	downloadHandler := kithttp.NewServer(NewDownloadEndpoint(svc), decodeDownloadRequest, encodeDownloadResponse, opts...)
	deleteFileHandler := kithttp.NewServer(NewDeleteFileEndpoint(svc), encodeDeleteFileRequest, encodeResponse, opts...)

	r.Handle("/api/v1/trips", createHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/trips", listHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}", readHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}/members", readMembersHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}", deleteHandler).Methods(http.MethodDelete)
	r.Handle("/api/v1/trips/{id}/storage", uploadHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/trips/{id}/storage", downloadHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}/storage", deleteFileHandler).Methods(http.MethodDelete)

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

func decodeUploadRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		fmt.Println(err)
		return nil, err
	}
	file, fileheader, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	buff := make([]byte, 512)
	if _, err = file.Read(buff); err != nil {
		return nil, err
	}
	filetype := http.DetectContentType(buff)
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	return UploadRequest{
		ID:             ID,
		Filename:       fileheader.Filename,
		Filesize:       fileheader.Size,
		AttachmentType: r.FormValue("attachmentType"),
		MimeType:       filetype,
		File:           file,
	}, nil
}

func encodeUploadResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(common.Errorer); ok && e.Error() != nil {
		common.EncodeErrorFactory(errToHttpCode)(ctx, e.Error(), w)
		return nil
	}
	resp, ok := response.(UploadResponse)
	if !ok {
		return common.ErrorEncodeInvalidResponse
	}
	resp.File.Close()
	return json.NewEncoder(w).Encode(response)
}

func decodeDownloadRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	path := r.URL.Query().Get("path")
	return DownloadRequest{
		ID: ID,
		Obj: storage.Object{
			Path: path,
			Name: filepath.Base(path),
		},
	}, nil
}

func encodeDownloadResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(common.Errorer); ok && e.Error() != nil {
		common.EncodeErrorFactory(errToHttpCode)(ctx, e.Error(), w)
		return nil
	}
	resp, ok := response.(DownloadResponse)
	if !ok {
		return common.ErrorEncodeInvalidResponse
	}

	//Set the headers
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", resp.Name))
	w.Header().Set("Content-Type", resp.Stat.MIMEType+";"+resp.Name)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", resp.Stat.Size))

	fmt.Println("copying")
	written, err := io.Copy(w, resp.File)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(written)

	resp.File.Close()
	return nil
}

func encodeDeleteFileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, common.ErrInvalidRequest
	}
	path := r.URL.Query().Get("path")
	return DeleteFileRequest{
		ID: ID,
		Obj: storage.Object{
			Path: path,
			Name: filepath.Base(path),
		},
	}, nil
}
