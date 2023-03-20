package tripssync

import (
	context "context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/travelreys/travelreys/pkg/common"
	"go.uber.org/zap"
)

// https://github.com/gorilla/websocket/tree/master/examples/chat

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketServer struct {
	svc    Service
	logger *zap.Logger
}

func NewWebsocketServer(svc Service, logger *zap.Logger) *WebsocketServer {
	return &WebsocketServer{svc: svc, logger: logger}
}

// HandleFunc upgrades the HTTP connection to the WebSocket protocol and then creates a ConnHandler.
func (srv *WebsocketServer) HandleFunc(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		srv.logger.Error("upgrader.Upgrade", zap.Error(err))
		return
	}
	defer ws.Close()

	h := ConnHandler{svc: srv.svc, ws: ws, logger: srv.logger}
	h.Run()
}

// ConnHandler handles a single Websocket connection. The proxy creates an instance of the ConnHandler type for each
// websocket connection. A ConnHandler acts as an intermediary between the websocket connection and the session.
// WebSocket connections support one concurrent reader and one concurrent writer.
// (https://pkg.go.dev/github.com/gorilla/websocket#hdr-Concurrency)
type ConnHandler struct {
	ID     string
	tripID string
	ws     *websocket.Conn

	svc      Service
	tobMsgCh <-chan Message
	doneCh   chan<- bool

	pongDeadline time.Time

	logger *zap.Logger
}

func (h *ConnHandler) SetPongDeadline(deadline time.Time) {
	h.pongDeadline = deadline
}

func (h *ConnHandler) Run() {
	h.logger.Info("new connection", zap.String("tripID", h.tripID))
	defer func() {
		h.logger.Info("closing connection", zap.String("id", h.ID))
		h.ReadMessage(NewMsgLeaveSession(h.ID, h.tripID))
		h.ws.Close()
	}()

	h.SetPongDeadline(time.Now().Add(pongWait))

	for {
		var msg Message
		if err := h.ws.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Error("h.ws.ReadJSON(&msg)", zap.Error(err))
			}
			return
		}
		if !isValidMessageType(msg.Op) {
			h.logger.Error("!isValidMessageType(msg.Op)", zap.Error(ErrInvalidOp))
			continue
		}
		if err := h.ReadMessage(msg); err != nil {
			h.logger.Error("h.HandleMessage(msg)", zap.Error(err))
			continue
		}
		// Close session if client leaves.
		if msg.Op == OpLeaveSession {
			return
		}
	}
}

func (h *ConnHandler) ReadMessage(msg Message) error {
	h.logger.Debug("recv msg", zap.String("op", msg.Op), zap.String("msg", common.FmtString(msg)))

	ctx := context.Background()

	switch msg.Op {
	case OpJoinSession:
		h.ID = msg.ConnID
		h.logger.Info("new client", zap.String("connID", msg.ConnID))
		tobMsgCh, doneCh, err := h.svc.SubscribeTOBUpdates(ctx, msg.TripID)
		if err != nil {
			return err
		}

		h.tobMsgCh = tobMsgCh
		h.tripID = msg.TripID
		h.doneCh = doneCh

		if _, err = h.svc.JoinSession(ctx, msg); err != nil {
			return err
		}
		go h.WriteMessage()
		return nil
	case OpPingSession:
		h.logger.Debug("pong")
		h.SetPongDeadline(time.Now().Add(pongWait))
		return nil
	case OpLeaveSession:
		h.doneCh <- true
		return h.svc.LeaveSession(ctx, msg)
	case OpUpdateTrip:
		return h.svc.UpdateTrip(ctx, msg)
	default:
		return nil
	}
}

func (h *ConnHandler) WriteMessage() {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
	}()

	for {
		select {
		case msg, ok := <-h.tobMsgCh:
			if !ok {
				return
			}
			h.logger.Debug("recv tob", zap.String("msg", common.FmtString(msg)))
			h.ws.WriteJSON(msg)
		case <-pingTicker.C:
			h.logger.Debug("ping")
			h.ws.WriteJSON(NewMsgPing(h.ID, h.tripID))
		}
	}
}
