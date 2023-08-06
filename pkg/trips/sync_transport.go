package trips

import (
	context "context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vmihailenco/msgpack/v5"
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

type ErrMessage struct {
	Err string `json:"error,omitempty"`
}

type WebsocketServer struct {
	svc    SyncService
	logger *zap.Logger
}

func NewWebsocketServer(svc SyncService, logger *zap.Logger) *WebsocketServer {
	return &WebsocketServer{svc: svc, logger: logger}
}

// HandleFunc upgrades the HTTP connection to the WebSocket
// protocol and then creates a ConnHandler.
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

// ConnHandler handles a single Websocket connection. The proxy creates an instance
// of the ConnHandler type for each websocket connection. A ConnHandler acts as
// an intermediary between the websocket connection and the session.
// WebSocket connections support one concurrent reader and one concurrent writer.
// (https://pkg.go.dev/github.com/gorilla/websocket#hdr-Concurrency)
type ConnHandler struct {
	ws *websocket.Conn

	connID   string
	tripID   string
	memberID string

	svc       SyncService
	dataMsgCh <-chan SyncMsgData
	doneCh    chan<- bool

	pongDeadline time.Time

	logger *zap.Logger
}

func (h *ConnHandler) SetPongDeadline(deadline time.Time) {
	h.pongDeadline = deadline
}

func (h *ConnHandler) Run() {
	h.logger.Info("new connection", zap.String("tripID", h.tripID))
	defer func() {
		h.logger.Info("closing connection", zap.String("id", h.connID))
		h.ReadControlMsg(MakeSyncMsgControlTopicLeave(h.connID, h.tripID, h.memberID))
		h.ws.Close()
	}()

	h.SetPongDeadline(time.Now().Add(pongWait))

	for {
		_, p, err := h.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				h.logger.Error("ReadMessage", zap.Error(err))
			}
			return
		}

		var syncMsg SyncMsg
		if err := msgpack.Unmarshal(p, &syncMsg); err != nil {
			h.logger.Error("msgpack.Unmarshal", zap.Error(err))
			return
		}

		if syncMsg.Type == SyncMsgTypeControl {
			var msg SyncMsgControl
			if err := msgpack.Unmarshal(p, &msg); err != nil {
				h.logger.Error("msgpack.Unmarshal", zap.Error(err))
				return
			}

			if err := h.ReadControlMsg(msg); err != nil {
				if err == ErrRBAC {
					b, _ := msgpack.Marshal(ErrMessage{Err: err.Error()})
					h.ws.WriteMessage(websocket.BinaryMessage, b)
					return
				}
				h.logger.Error("ReadControlMsg", zap.Error(err))
				continue
			}

			// Close session if client leaves.
			if msg.Topic == SyncMsgControlTopicLeave {
				return
			}
			continue
		}

		var msg SyncMsgData
		if err := msgpack.Unmarshal(p, &msg); err != nil {
			h.logger.Error("msgpack.Unmarshal", zap.Error(err))
			return
		}

		if err := h.ReadDataMessage(msg); err != nil {
			h.logger.Error("h.ReadDataMessage", zap.Error(err))
			continue
		}

	}
}

func (h *ConnHandler) ReadDataMessage(msg SyncMsgData) error {
	h.logger.Debug("recv data msg", zap.String("op", msg.Topic))
	return h.svc.UpdateTrip(context.Background(), msg)
}

func (h *ConnHandler) ReadControlMsg(msg SyncMsgControl) error {
	h.logger.Debug("recv control msg", zap.String("topic", msg.Topic))

	ctx := context.Background()

	switch msg.Topic {
	case SyncMsgControlTopicJoin:
		h.connID = msg.ConnID
		h.logger.Info("new client", zap.String("connID", msg.ConnID))
		dataMsgCh, doneCh, err := h.svc.SubSyncMsgDataResp(ctx, msg.TripID)
		if err != nil {
			return err
		}

		h.dataMsgCh = dataMsgCh
		h.doneCh = doneCh
		h.tripID = msg.TripID
		h.memberID = msg.MemberID
		if err = h.svc.Join(ctx, msg); err != nil {
			return err
		}
		go h.WriteMessage()
		return nil
	case SyncMsgControlTopicPing:
		h.logger.Debug("pong")
		h.SetPongDeadline(time.Now().Add(pongWait))
		msg.MemberID = h.memberID
		h.svc.Ping(ctx, msg)
		return nil
	case SyncMsgControlTopicLeave:
		h.doneCh <- true
		return h.svc.Leave(ctx, msg)
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
		case msg, ok := <-h.dataMsgCh:
			if !ok {
				return
			}
			msg.ConnID = h.connID
			h.logger.Debug("recv tob", zap.String("op", msg.Topic))
			b, _ := msgpack.Marshal(msg)
			h.ws.WriteMessage(websocket.BinaryMessage, b)
		case <-pingTicker.C:
			h.logger.Debug("ping")
			b, _ := msgpack.Marshal(MakeSyncMsgControlTopicPing(h.connID, h.tripID, h.memberID))
			h.ws.WriteMessage(websocket.BinaryMessage, b)
		}
	}
}
