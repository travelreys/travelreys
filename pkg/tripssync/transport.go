package tripssync

import (
	context "context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var (
	ErrInvalidSyncOp     = errors.New("invalid-collab-op")
	ErrInvalidSyncOpData = errors.New("invalid-collab-op-data")
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

/**********
 * Server *
 *********/

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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

	go srv.PingPong(ws)

	connHandler := ConnHandler{
		proxy:  srv.proxy,
		connID: uuid.New().String(),
		ws:     ws,
		logger: srv.logger,
	}
	connHandler.Run()
}

func (srv *Server) PingPong(ws *websocket.Conn) {
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()
	for range pingTicker.C {
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			return
		}
	}
}

// Connection Handler

type ConnHandler struct {
	proxy     Proxy
	syncMsgCh chan SyncMessage
	ws        *websocket.Conn

	connID     string
	tripPlanID string

	logger *zap.Logger
}

func (handler *ConnHandler) Run() {
	defer func() {
		handler.HandleSyncMessage(SyncMessage{
			OpType:     SyncOpLeaveSession,
			ID:         handler.connID,
			TripPlanID: handler.tripPlanID,
		})
		handler.ws.Close()
		close(handler.syncMsgCh)
	}()

	for {
		var msg SyncMessage
		err := handler.ws.ReadJSON(&msg)
		if err != nil {
			handler.logger.Error("read:", zap.Error(err))
			return
		}
		if !isValidSyncOpType(msg.OpType) {
			continue
		}
		if err := handler.HandleSyncMessage(msg); err != nil {
			handler.logger.Error("handle:", zap.Error(err))
			continue
		}
		// Close session if client leaves.
		if msg.OpType == SyncOpLeaveSession {
			return
		}

		// handler.ws.WriteJSON(resp)
	}
}

func (handler *ConnHandler) HandleSyncMessage(msg SyncMessage) error {
	ctx := context.Background()

	msg.ID = handler.connID

	switch msg.OpType {
	case SyncOpJoinSession:
		msgCh, err := handler.proxy.SubscribeTOBUpdates(context.Background(), msg.TripPlanID)
		if err != nil {
			return err
		}
		handler.syncMsgCh = msgCh
		handler.tripPlanID = msg.TripPlanID
		_, err = handler.proxy.JoinSession(ctx, msg.TripPlanID, msg)
		if err == nil {
			go handler.HandleProxy()
		}
		return err

	case SyncOpLeaveSession:
		return handler.proxy.LeaveSession(ctx, msg.TripPlanID, msg)

	case SyncOpUpdateTrip:
		return handler.proxy.UpdateTripPlan(ctx, msg.TripPlanID, msg)

	default:
		return nil
	}
}

func (handler *ConnHandler) HandleProxy() {
	for msg := range handler.syncMsgCh {
		fmt.Println(msg)
		return
	}
}
