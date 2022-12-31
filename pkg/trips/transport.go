package trips

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/awhdesmond/tiinyplanet/pkg/common"
	"github.com/awhdesmond/tiinyplanet/pkg/reqctx"
	"github.com/awhdesmond/tiinyplanet/pkg/utils"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var (
	ErrInvalidCollabOp     = errors.New("invalid-collab-op")
	ErrInvalidCollabOpData = errors.New("invalid-collab-op-data")
)

// Trips Transport

const (
	URLPathVarID = "id"
)

var (
	encodeErrFn = utils.EncodeErrorFactory(ErrorToHTTPCode)

	opts = []kithttp.ServerOption{
		kithttp.ServerBefore(reqctx.MakeContextFromHTTPRequest),
		kithttp.ServerErrorEncoder(encodeErrFn),
	}
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(common.Errorer); ok && e.Error() != nil {
		encodeErrFn(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func MakeHandler(svc Service) http.Handler {
	r := mux.NewRouter()
	createTripPlanHandler := kithttp.NewServer(NewCreateTripPlanEndpoint(svc), decodeCreateTripPlanRequest, encodeResponse, opts...)
	listTripPlansHandler := kithttp.NewServer(NewListTripPlansEndpoint(svc), decodeListTripPlansRequest, encodeResponse, opts...)
	readTripPlanHandler := kithttp.NewServer(NewReadTripPlanEndpoint(svc), decodeReadTripPlanRequest, encodeResponse, opts...)
	deleteTripPlanHandler := kithttp.NewServer(NewDeleteTripPlanEndpoint(svc), decodeDeleteTripPlanRequest, encodeResponse, opts...)

	r.Handle("/api/v1/trips", createTripPlanHandler).Methods(http.MethodPost)
	r.Handle("/api/v1/trips", listTripPlansHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}", readTripPlanHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/trips/{id}", deleteTripPlanHandler).Methods(http.MethodDelete)

	return r
}

func decodeCreateTripPlanRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := CreateTripPlanRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, utils.ErrInvalidJSONBody
	}
	return req, nil
}
func decodeReadTripPlanRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, utils.ErrInvalidRequest
	}
	return ReadTripPlanRequest{ID}, nil
}
func decodeListTripPlansRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return ListTripPlansRequest{}, nil

}
func decodeDeleteTripPlanRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	ID, ok := vars[URLPathVarID]
	if !ok {
		return nil, utils.ErrInvalidRequest
	}
	return DeleteTripPlanRequest{ID}, nil
}

// Collab Transport

type CollabServer struct {
	UnimplementedCollabServiceServer

	collabService CollabService
	logger        *zap.Logger
}

func MakeCollabServer(cs CollabService, logger *zap.Logger) *CollabServer {
	return &CollabServer{collabService: cs, logger: logger}
}

func (srv *CollabServer) Command(tsrv CollabService_CommandServer) error {
	ctx := tsrv.Context()

	// srv.collabService.collabStore.SubscribeCollabOpMessages()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		req, err := tsrv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			continue
		}

		var resp CollabResponse

		opType := req.GetType()
		if !isValidCollabOpType(opType) {
			resp.Error = ErrInvalidCollabOp.Error()
		} else {
			result, err := srv.HandleCollabRequest(req)
			resp = CollabResponse{
				Error:  err.Error(),
				Type:   opType,
				Result: result,
			}
		}
		if err := tsrv.Send(&resp); err != nil {
			log.Printf("send error %v", err)
		}
	}
}

func (srv *CollabServer) HandleCollabRequest(req *CollabRequest) ([]byte, error) {
	var msg CollabOpMessage
	if err := json.Unmarshal(req.GetData(), &msg); err != nil {
		return nil, ErrInvalidCollabOpData
	}

	opType := req.GetType()
	switch opType {
	case CollabOpJoinSession:
		session, err := srv.collabService.JoinSession(reqctx.Context{}, msg.TripPlanID, msg)
		result := common.GenericJSON{"session": session, "error": err}
		resBytes, _ := json.Marshal(result)
		return resBytes, nil
	case CollabOpLeaveSession:
		err := srv.collabService.LeaveSession(reqctx.Context{}, msg.TripPlanID, msg)
		result := common.GenericJSON{"error": err}
		resBytes, _ := json.Marshal(result)
		return resBytes, nil
	default:
		return nil, ErrInvalidCollabOp
	}
}
