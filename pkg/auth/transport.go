package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/travelreys/travelreys/pkg/common"
	"github.com/travelreys/travelreys/pkg/reqctx"
	"github.com/travelreys/travelreys/pkg/storage"
)

const (
	AccessCookieName = "_travelreysAuth"
	mediaValidCookie = "_travelreys_media"
	cdnCookie        = "Cloud-CDN-Cookie"
)

func errToHttpCode(err error) int {
	notFoundErrors := []error{ErrUserNotFound}
	appErrors := []error{
		ErrProviderNotSupported,
		ErrProviderGoogleError,
		ErrProviderOTPEmailNotFound,
		ErrProviderOTPEmailExists,
		ErrProviderOTPNotSet,
		ErrProviderOTPInvalidEmail,
		ErrProviderOTPInvalidPw,
		ErrProviderOTPInvalidSig,
	}
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

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	if e, ok := resp.(common.Errorer); ok && e.Error() != nil {
		common.EncodeErrorFactory(errToHttpCode)(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(resp)
}

func MakeHandler(svc Service) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(reqctx.ContextWithClientInfo),
		kithttp.ServerErrorEncoder(common.EncodeErrorFactory(errToHttpCode)),
	}

	loginHandler := kithttp.NewServer(NewLoginEndpoint(svc), decodeLoginRequest, encodeLoginResponse, opts...)
	logoutHandler := kithttp.NewServer(NewLogoutEndpoint(svc), decodeLogoutRequest, encodeLogoutResponse, opts...)
	magicLinkHandler := kithttp.NewServer(NewMagicLinkEndpoint(svc), decodeMagicLinkRequest, encodeResponse, opts...)
	readUserHandler := kithttp.NewServer(NewReadEndpoint(svc), decodeReadRequest, encodeResponse, opts...)
	listUsersHandler := kithttp.NewServer(NewListEndpoint(svc), decodeListRequest, encodeResponse, opts...)
	updateUserHandler := kithttp.NewServer(NewUpdateEndpoint(svc), decodeUpdateRequest, encodeResponse, opts...)
	deleteUserHandler := kithttp.NewServer(NewDeleteEndpoint(svc), decodeDeleteRequest, encodeDeleteResponse, opts...)
	uploadAvatarPresignedURLHandler := kithttp.NewServer(
		NewUploadAvatarPresignedURLEndpoint(svc),
		decodeUploadAvatarPresignedURLRequest,
		encodeResponse,
		opts...,
	)
	generateMediaPresignedCookieHandler := kithttp.NewServer(
		NewGenerateMediaPresignedCookieEndpoint(svc),
		decodeGenerateMediaPresignedCookieRequest,
		encodeGenerateMediaPresignedCookieResponse,
		opts...,
	)

	r.Handle("/api/v1/auth/login", loginHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/auth/logout", logoutHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/auth/magic-link", magicLinkHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/auth/users", listUsersHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/auth/users/{id}", readUserHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/auth/users/{id}", updateUserHandler).Methods(http.MethodPut)
	r.Handle("/api/v1/auth/users/{id}", deleteUserHandler).Methods(http.MethodDelete)
	r.Handle("/api/v1/auth/users/{id}/avatar/upload/pre-signed", uploadAvatarPresignedURLHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/auth/media/pre-signed-cookie", generateMediaPresignedCookieHandler).Methods(http.MethodGet)

	return r
}

func decodeMagicLinkRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := MagicLinkRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	return req, nil
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	return req, nil
}

func encodeLoginResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(common.Errorer); ok && e.Error() != nil {
		common.EncodeErrorFactory(errToHttpCode)(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp, ok := response.(LoginResponse)
	if !ok {
		return common.ErrorEncodeInvalidResponse
	}
	http.SetCookie(w, resp.Cookie)
	return json.NewEncoder(w).Encode(response)
}

func decodeLogoutRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return LogoutRequest{}, nil
}

func encodeLogoutResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	accessCookie := &http.Cookie{
		Name:     AccessCookieName,
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now(),
	}
	http.SetCookie(w, accessCookie)
	mediaCookie := &http.Cookie{
		Name:     mediaValidCookie,
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now(),
	}
	http.SetCookie(w, mediaCookie)
	cdnCookie := &http.Cookie{
		Name:     cdnCookie,
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now(),
	}
	http.SetCookie(w, cdnCookie)
	return nil
}

func decodeReadRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	req := ReadRequest{
		ID: vars[bsonKeyID],
	}
	return req, nil
}

func decodeUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	req := UpdateRequest{
		ID: vars[bsonKeyID],
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, common.ErrInvalidJSONBody
	}
	return req, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := ListRequest{ListFilter{}}
	if r.URL.Query().Has(bsonKeyEmail) {
		req.FF.Email = common.StringPtr(r.URL.Query().Get(bsonKeyEmail))
	}
	return req, nil
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	req := DeleteRequest{
		ID: vars[bsonKeyID],
	}
	return req, nil
}

func encodeDeleteResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	accessCookie := &http.Cookie{
		Name:     AccessCookieName,
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now(),
	}
	http.SetCookie(w, accessCookie)
	mediaCookie := &http.Cookie{
		Name:     mediaValidCookie,
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now(),
	}
	http.SetCookie(w, mediaCookie)
	cdnCookie := &http.Cookie{
		Name:     cdnCookie,
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now(),
	}
	http.SetCookie(w, cdnCookie)
	return nil
}

func decodeUploadAvatarPresignedURLRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	req := UploadAvatarPresignedURLRequest{
		ID: vars[bsonKeyID],
	}
	return req, nil
}

func decodeGenerateMediaPresignedCookieRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return GenerateMediaPresignedCookieRequest{}, nil
}

func encodeGenerateMediaPresignedCookieResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(common.Errorer); ok && e.Error() != nil {
		common.EncodeErrorFactory(errToHttpCode)(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp, ok := response.(GenerateMediaPresignedCookieResponse)
	if !ok {
		return common.ErrorEncodeInvalidResponse
	}

	if resp.Cookie != nil {
		http.SetCookie(w, resp.Cookie)
		expCookie := &http.Cookie{
			Name:   mediaValidCookie,
			Value:  "1",
			Path:   "/",
			MaxAge: int(storage.DefaultPresignedCookieRefreshDuration.Seconds()),
		}
		http.SetCookie(w, expCookie)
	}
	return json.NewEncoder(w).Encode(response)
}
