package tripssync

import (
	context "context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var (
	ErrInvalidCollabOp     = errors.New("invalid-collab-op")
	ErrInvalidCollabOpData = errors.New("invalid-collab-op-data")
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
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

var upgrader = websocket.Upgrader{} // use default options

type Server struct {
	proxy  Proxy
	logger *zap.Logger
}

func MakeServer(pxy Proxy, logger *zap.Logger) *Server {
	return &Server{proxy: pxy, logger: logger}
}

func (srv *Server) HandleFunc(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		srv.logger.Error("upgrade", zap.Error(err))
		return
	}
	defer ws.Close()

	pingTicker := time.NewTicker(pingPeriod)
	go func() {
		for range pingTicker.C {
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}()

	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		var msg CollabOpMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("read:", err)
			continue
		}

		if !isValidCollabOpType(msg.OpType) {
			ws.WriteJSON(GenericSessionResponse{err})
			continue
		}

		resp, err := srv.HandleCollabMessage(msg)
		if err != nil {
			log.Println("handle:", err)
			continue
		}
		ws.WriteJSON(resp)
	}
}

func (srv *Server) HandleCollabMessage(msg CollabOpMessage) (interface{}, error) {
	ctx := context.Background()

	switch msg.OpType {
	case CollabOpJoinSession:
		session, err := srv.proxy.JoinSession(ctx, msg.TripPlanID, msg)
		return JoinSessionResponse{session, err}, nil
	case CollabOpLeaveSession:
		err := srv.proxy.LeaveSession(ctx, msg.TripPlanID, msg)
		return GenericSessionResponse{err}, nil
	case CollabOpUpdateTrip:
		err := srv.proxy.UpdateTripPlan(ctx, msg.TripPlanID, msg)
		return GenericSessionResponse{err}, nil
	default:
		return nil, ErrInvalidCollabOp
	}
}
