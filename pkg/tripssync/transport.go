package tripssync

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"

	"go.uber.org/zap"
)

var (
	ErrInvalidCollabOp     = errors.New("invalid-collab-op")
	ErrInvalidCollabOpData = errors.New("invalid-collab-op-data")
)

/*************
 * Responses *
 *************/

type JoinSessionResponse struct {
	Session CollabSession `json:"session"`
	Err     error         `json:"error,omitempty"`
}

type GenericSessionResponse struct {
	Err error `json:"error,omitempty"`
}

/**********
 * Server *
 *********/

type Server struct {
	UnimplementedCollabServiceServer

	proxy  Proxy
	logger *zap.Logger
}

func MakeServer(pxy Proxy, logger *zap.Logger) *Server {
	return &Server{proxy: pxy, logger: logger}
}

func (srv *Server) sendResp(resp *CollabResponse, tsrv CollabService_CommandServer) {
	if err := tsrv.Send(resp); err != nil {
		log.Printf("send error %v", err)
	}
}

func (srv *Server) Command(tsrv CollabService_CommandServer) error {
	ctx := tsrv.Context()
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
			srv.sendResp(&resp, tsrv)
			continue
		}

		result, err := srv.HandleCollabRequest(req)
		resp = CollabResponse{Error: err.Error(), Type: opType, Result: result}
		srv.sendResp(&resp, tsrv)
	}
}

func (srv *Server) HandleCollabRequest(req *CollabRequest) ([]byte, error) {
	var msg CollabOpMessage
	if err := json.Unmarshal(req.GetData(), &msg); err != nil {
		return nil, ErrInvalidCollabOpData
	}

	ctx := context.Background()

	opType := req.GetType()
	switch opType {
	case CollabOpJoinSession:
		session, err := srv.proxy.JoinSession(ctx, msg.TripPlanID, msg)
		resp := JoinSessionResponse{session, err}
		respBytes, _ := json.Marshal(resp)
		return respBytes, nil
	case CollabOpLeaveSession:
		err := srv.proxy.LeaveSession(ctx, msg.TripPlanID, msg)
		resp := GenericSessionResponse{err}
		respBytes, _ := json.Marshal(resp)
		return respBytes, nil
	case CollabOpUpdateTrip:
		err := srv.proxy.UpdateTripPlan(ctx, msg.TripPlanID, msg)
		resp := GenericSessionResponse{err}
		respBytes, _ := json.Marshal(resp)
		return respBytes, nil
	default:
		return nil, ErrInvalidCollabOp
	}
}
