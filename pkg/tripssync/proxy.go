package tripssync

import (
	context "context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tiinyplanet/tiinyplanet/pkg/trips"
	"go.uber.org/zap"
)

// Proxy proxies clients' updates to the backend sync infrastructure.
type Proxy interface {
	JoinSession(ctx context.Context, planID string, msg SyncMessage) (SyncSession, error)
	LeaveSession(ctx context.Context, planID string, msg SyncMessage) error
	ReadTripPlan(ctx context.Context, planID string, msg SyncMessage) (trips.TripPlan, error)
	UpdateTripPlan(ctx context.Context, planID string, msg SyncMessage) error
	SubscribeTOBUpdates(ctx context.Context, planID string) (<-chan SyncMessage, chan<- bool, error)
}

type proxy struct {
	sesnStore SessionStore
	sms       SyncMessageStore
	tms       TOBMessageStore
	tripStore trips.Store
}

func NewProxy(
	sesnStore SessionStore,
	sms SyncMessageStore,
	tms TOBMessageStore,
	tripStore trips.Store,
) (Proxy, error) {
	return &proxy{sesnStore, sms, tms, tripStore}, nil
}

// Session

func (p *proxy) JoinSession(ctx context.Context, planID string, msg SyncMessage) (SyncSession, error) {
	conn := SyncConnection{
		PlanID:       planID,
		ConnectionID: msg.ID,
		Member:       msg.SyncDataJoinSession.TripMember,
	}
	err := p.sesnStore.AddConnToSession(ctx, conn)
	if err != nil {
		return SyncSession{}, err
	}

	p.sms.Publish(planID, msg)
	return p.sesnStore.Read(ctx, planID)
}

func (p *proxy) LeaveSession(ctx context.Context, planID string, msg SyncMessage) error {
	conn := SyncConnection{
		PlanID:       planID,
		ConnectionID: msg.ID,
	}
	p.sesnStore.RemoveConnFromSession(ctx, conn)
	p.sms.Publish(planID, msg)
	return nil
}

// Plans

func (p *proxy) ReadTripPlan(ctx context.Context, planID string, msg SyncMessage) (trips.TripPlan, error) {
	return p.tripStore.ReadTripPlan(ctx, planID)
}

// Sync Messages

func (p *proxy) UpdateTripPlan(ctx context.Context, planID string, msg SyncMessage) error {
	return p.sms.Publish(planID, msg)
}

// TOB Updates

func (p *proxy) SubscribeTOBUpdates(ctx context.Context, planID string) (<-chan SyncMessage, chan<- bool, error) {
	return p.tms.Subscribe(planID)
}

/****************
 * Proxy Server *
 ****************/

var (
	ErrInvalidSyncOp     = errors.New("invalid-sync-op")
	ErrInvalidSyncOpData = errors.New("invalid-sync-op-data")
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

// ProxyServer handles multiple incoming websocket connections from clients
// by creating a connection handler for each connection.
type ProxyServer struct {
	proxy  Proxy
	logger *zap.Logger
}

func MakeProxyServer(pxy Proxy, logger *zap.Logger) *ProxyServer {
	return &ProxyServer{proxy: pxy, logger: logger}
}

func (srv *ProxyServer) HandleFunc(w http.ResponseWriter, r *http.Request) {
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

func (srv *ProxyServer) PingPong(ws *websocket.Conn) {
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
	proxy      Proxy
	ws         *websocket.Conn
	connID     string
	tripPlanID string
	tobMsgCh   <-chan SyncMessage
	doneCh     chan<- bool
	logger     *zap.Logger
}

func (h *ConnHandler) Run() {
	h.logger.Info("new connection", zap.String("connID", h.connID))
	defer func() {
		h.logger.Info("closing connection", zap.String("connID", h.connID))
		h.HandleSyncMessage(NewSyncMessageLeaveSession(h.connID, h.tripPlanID))
		h.ws.Close()
	}()

	for {
		var msg SyncMessage
		err := h.ws.ReadJSON(&msg)
		if err != nil {
			h.logger.Error("read:", zap.Error(err))
			return
		}
		if !isValidSyncMessageType(msg.OpType) {
			continue
		}
		if err := h.HandleSyncMessage(msg); err != nil {
			h.logger.Error("handle:", zap.Error(err))
			continue
		}
		// Close session if client leaves.
		if msg.OpType == SyncOpLeaveSession {
			return
		}

	}
}

func (h *ConnHandler) HandleSyncMessage(msg SyncMessage) error {
	ctx := context.Background()
	msg.ID = h.connID

	h.logger.Debug("recv msg", zap.String("msg", fmt.Sprintf("%+v", msg)))

	switch msg.OpType {
	case SyncOpJoinSession:
		tobMsgCh, done, err := h.proxy.SubscribeTOBUpdates(context.Background(), msg.TripPlanID)
		if err != nil {
			return err
		}

		h.tobMsgCh = tobMsgCh
		h.tripPlanID = msg.TripPlanID
		h.doneCh = done

		_, err = h.proxy.JoinSession(ctx, msg.TripPlanID, msg)
		if err == nil {
			go h.HandleProxy()
		}
		return err

	case SyncOpLeaveSession:
		h.doneCh <- true
		return h.proxy.LeaveSession(ctx, msg.TripPlanID, msg)

	case SyncOpUpdateTrip:
		return h.proxy.UpdateTripPlan(ctx, msg.TripPlanID, msg)

	default:
		return nil
	}
}

func (h *ConnHandler) HandleProxy() {
	for msg := range h.tobMsgCh {
		h.logger.Debug("recv tob", zap.String("msg", fmt.Sprintf("%+v", msg)))
		h.ws.WriteJSON(msg)
	}
}
