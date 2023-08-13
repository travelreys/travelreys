package trips

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

	svc        SyncService
	ctrlMsgCh  <-chan SyncMsgBroadcast
	ctrlDoneCh chan<- bool
	dataMsgCh  <-chan SyncMsgTOB
	dataDoneCh chan<- bool

	pongDeadline time.Time

	logger *zap.Logger
}

func (h *ConnHandler) SetPongDeadline(deadline time.Time) {
	h.pongDeadline = deadline
}

func (h *ConnHandler) Run() {
	h.logger.Info("new connection")
	defer func() {
		h.logger.Info("closing connection", zap.String("id", h.connID))
		msg := MakeSyncMsgTOBTopicLeave(h.connID, h.tripID, h.memberID)
		h.ProcessDataMessage(&msg)
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
		if err := common.MsgpackUnmarshal(p, &syncMsg); err != nil {
			h.logger.Error("common.MsgpackUnmarshal", zap.Error(err))
			continue
		}

		switch syncMsg.Type {
		case SyncMsgTypeBroadcast:
			ctrlmsg, err := h.ParseControlMsg(p)
			if err != nil {
				h.logger.Warn("ParseControlMsg", zap.Error(err))
				continue
			}
			if err := h.ProcessControlMsg(ctrlmsg); err != nil {
				h.logger.Error("ProcessControlMsg", zap.Error(err))
				continue
			}
			// Close session if client leaves.
			if ctrlmsg.Topic == SyncMsgTOBTopicLeave {
				return
			}
			continue
		case SyncMsgTypeTOB:
			var msg SyncMsgTOB
			if err := common.MsgpackUnmarshal(p, &msg); err != nil {
				h.logger.Error("common.MsgpackUnmarshal", zap.Error(err))
				continue
			}

			if err := h.ProcessDataMessage(&msg); err != nil {
				if err == ErrRBAC {
					b, _ := common.MsgpackMarshal(ErrMessage{Err: err.Error()})
					h.ws.WriteMessage(websocket.BinaryMessage, b)
					return
				}
				h.logger.Error("h.ProcessDataMessage", zap.Error(err))
				continue
			}
		}
	}
}

func (h *ConnHandler) ParseControlMsg(p []byte) (*SyncMsgBroadcast, error) {
	var msg SyncMsgBroadcast
	if err := common.MsgpackUnmarshal(p, &msg); err != nil {
		h.logger.Error("common.MsgpackUnmarshal", zap.Error(err))
		return nil, err
	}
	return &msg, nil
}

func (h *ConnHandler) ProcessControlMsg(msg *SyncMsgBroadcast) error {
	h.logger.Debug("recv control msg", zap.String("topic", msg.Topic))

	ctx := context.Background()

	switch msg.Topic {
	case SyncMsgBroadcastTopicPing:
		h.logger.Debug("pong")
		h.SetPongDeadline(time.Now().Add(pongWait))
		msg.MemberID = h.memberID
		h.svc.Ping(ctx, msg)
		return nil

	default:
		return nil
	}
}

func (h *ConnHandler) ProcessDataMessage(msg *SyncMsgTOB) error {
	h.logger.Debug("recv data msg", zap.String("op", msg.Topic))

	ctx := context.Background()
	switch msg.Topic {
	case SyncMsgTOBTopicJoin:
		h.connID = msg.ConnID
		h.logger.Info("new client", zap.String("connID", msg.ConnID))

		ctrlMsgCh, ctrlDoneCh, err := h.svc.SubSyncMsgBroadcastResp(ctx, msg.TripID)
		if err != nil {
			return err
		}
		h.ctrlMsgCh = ctrlMsgCh
		h.ctrlDoneCh = ctrlDoneCh

		dataMsgCh, dataDoneCh, err := h.svc.SubSyncMsgTOBResp(ctx, msg.TripID)
		if err != nil {
			return err
		}
		h.dataMsgCh = dataMsgCh
		h.dataDoneCh = dataDoneCh
		h.tripID = msg.TripID
		h.memberID = msg.MemberID
		if err = h.svc.Join(ctx, msg); err != nil {
			return err
		}
		go h.WriteMessage()
		return nil
	case SyncMsgTOBTopicLeave:
		h.ctrlDoneCh <- true
		h.dataDoneCh <- true
		return h.svc.Leave(ctx, msg)
	}

	return h.svc.Update(context.Background(), msg)
}

func (h *ConnHandler) WriteMessage() {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
	}()

	for {
		select {
		case msg, ok := <-h.ctrlMsgCh:
			if !ok {
				return
			}
			msg.ConnID = h.connID
			h.logger.Debug("recv control", zap.String("op", msg.Topic))
			b, _ := common.MsgpackMarshal(msg)
			h.ws.WriteMessage(websocket.BinaryMessage, b)
		case msg, ok := <-h.dataMsgCh:
			if !ok {
				return
			}
			msg.ConnID = h.connID
			h.logger.Debug("recv tob", zap.String("op", msg.Topic))
			b, _ := common.MsgpackMarshal(msg)
			h.ws.WriteMessage(websocket.BinaryMessage, b)
		case <-pingTicker.C:
			h.logger.Debug("ping")
			b, _ := common.MsgpackMarshal(
				MakeSyncMsgBroadcastTopicPing(h.connID, h.tripID, h.memberID),
			)
			h.ws.WriteMessage(websocket.BinaryMessage, b)
		}
	}
}
