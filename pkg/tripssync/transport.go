package tripssync

import (
	context "context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tiinyplanet/tiinyplanet/pkg/common"
	"go.uber.org/zap"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
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

func (srv *WebsocketServer) HandleFunc(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		srv.logger.Error("upgrader.Upgrade", zap.Error(err))
		return
	}
	defer ws.Close()

	// go srv.PingPong(ws)
	connID := uuid.New().String()
	h := ConnHandler{ID: connID, svc: srv.svc, ws: ws, logger: srv.logger}
	h.Run()
}

func (srv *WebsocketServer) PingPong(ws *websocket.Conn) {
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		srv.logger.Debug("pingpong read")
		return nil
	})

	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()
	for range pingTicker.C {
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		srv.logger.Debug("pingpong write deadline")
		if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			srv.logger.Debug("pingpong write deadline error")
			srv.logger.Error("ping error", zap.Error(err))
			return
		}
	}
}

// ConnHandler handles a single Websocket connection
type ConnHandler struct {
	ID       string
	tripID   string
	ws       *websocket.Conn
	svc      Service
	tobMsgCh <-chan Message
	doneCh   chan<- bool
	logger   *zap.Logger
}

func (h *ConnHandler) Run() {
	h.logger.Info("new connection", zap.String("id", h.ID))
	defer func() {
		h.logger.Info("closing connection", zap.String("id", h.ID))
		h.HandleMessage(NewMsgLeaveSession(h.ID, h.tripID))
		h.ws.Close()
	}()

	for {
		var msg Message
		if err := h.ws.ReadJSON(&msg); err != nil {
			h.logger.Error("h.ws.ReadJSON(&msg)", zap.Error(err))
			return
		}
		if !isValidMessageType(msg.Op) {
			h.logger.Error("!isValidMessageType(msg.Op)", zap.Error(ErrInvalidOp))
			continue
		}
		if err := h.HandleMessage(msg); err != nil {
			h.logger.Error("h.HandleMessage(msg)", zap.Error(err))
			continue
		}
		// Close session if client leaves.
		if msg.Op == OpLeaveSession {
			return
		}

	}
}

func (h *ConnHandler) HandleMessage(msg Message) error {
	h.logger.Info("recv msg", zap.String("msg", common.FmtString(msg)))

	msg.ConnID = h.ID
	ctx := context.Background()

	switch msg.Op {
	case OpJoinSession:
		tobMsgCh, doneCh, err := h.svc.SubscribeTOBUpdates(ctx, msg.TripID)
		if err != nil {
			return err
		}

		h.tobMsgCh = tobMsgCh
		h.tripID = msg.TripID
		h.doneCh = doneCh

		if _, err = h.svc.JoinSession(ctx, msg.TripID, msg); err != nil {
			return err
		}
		go func() {
			for msg := range h.tobMsgCh {
				h.logger.Debug("recv tob", zap.String("msg", common.FmtString(msg)))
				h.ws.WriteJSON(msg)
			}
		}()
		return nil
	case OpLeaveSession:
		h.doneCh <- true
		return h.svc.LeaveSession(ctx, msg.TripID, msg)
	case OpUpdateTrip:
		return h.svc.UpdateTrip(ctx, msg.TripID, msg)
	default:
		return nil
	}
}
