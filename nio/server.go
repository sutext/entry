package nio

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/cloudwego/netpoll"
	"sutext.github.io/entry/logger"
	"sutext.github.io/entry/packet"
	"sutext.github.io/entry/safe"
)

type DataHandler func(uid string, p *packet.DataPacket) (*packet.DataPacket, error)
type AuthHandler func(p *packet.Identity) error
type NIOServer struct {
	port        string
	conns       *safe.Map[string, netpoll.Connection]
	groups      *safe.Map[string, []string]
	logger      *slog.Logger
	eventLoop   netpoll.EventLoop
	dataHandler DataHandler
	authHandler AuthHandler
}

func New(port string) *NIOServer {
	s := &NIOServer{
		port:   port,
		conns:  safe.NewMap(map[string]netpoll.Connection{}),
		groups: safe.NewMap(map[string][]string{}),
		logger: logger.New(logger.LevelDebug, logger.FormatJSON),
		authHandler: func(p *packet.Identity) error {
			return fmt.Errorf("login handler not set")
		},
		dataHandler: func(uid string, p *packet.DataPacket) (*packet.DataPacket, error) {
			return nil, fmt.Errorf("data handler not set")
		},
	}
	return s
}
func (s *NIOServer) SetLogger(level logger.Level, format logger.Format) {
	s.logger = logger.New(level, format)
}
func (s *NIOServer) HandleAuth(handler AuthHandler) {
	s.authHandler = handler
}
func (s *NIOServer) HandleData(handler DataHandler) {
	s.dataHandler = handler
}
func (s *NIOServer) Start() error {
	ln, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return err
	}
	eventLoop, err := netpoll.NewEventLoop(s.onRequest,
		netpoll.WithOnPrepare(s.onPrepare),
		netpoll.WithOnConnect(s.onConnect),
		netpoll.WithOnDisconnect(s.onDisconnect),
	)
	if err != nil {
		return err
	}
	s.eventLoop = eventLoop
	s.logger.Info("starting nio server", slog.String("port", s.port))
	err = eventLoop.Serve(ln)
	if err != nil {
		return err
	}
	return nil
}
func (s *NIOServer) Shutdown(ctx context.Context) error {
	if s.eventLoop != nil {
		return s.eventLoop.Shutdown(ctx)
	}
	return nil
}
func (s *NIOServer) close(clientID string, code packet.CloseCode) {
	if conn, ok := s.conns.Get(clientID); ok {
		packet.WriteTo(conn, packet.NewClose(code))
		conn.Close()
		s.conns.Delete(clientID)
	}
}
func (s *NIOServer) sendPacket(clientID string, p packet.Packet) error {
	if conn, ok := s.conns.Get(clientID); ok {
		return packet.WriteTo(conn, p)
	}
	return fmt.Errorf("conn not found")
}
func (s *NIOServer) handlePacket(id *packet.Identity, p packet.Packet) {
	switch p.Type() {
	case packet.DATA:
		p := p.(*packet.DataPacket)
		res, err := s.dataHandler(id.UserID, p)
		if err != nil {
			s.logger.Error("data handler failed", "error", err)
			return
		}
		if res != nil {
			s.sendPacket(id.ClientID, res)
		}
	case packet.PING:
		s.logger.Debug("receive ping", "client_id", id.ClientID)
		s.sendPacket(id.ClientID, packet.NewPong())
	case packet.CONNECT:
		break
	case packet.PONG:
		break
	case packet.CLOSE:
	default:
		break
	}
}

func (s *NIOServer) onRequest(ctx context.Context, conn netpoll.Connection) error {
	id := ctx.Value(identityKey).(*packet.Identity)
	pkt, err := packet.ReadFrom(conn)
	if err != nil {
		s.logger.Error("read packet failed", "error", err, "client_id", id.ClientID)
		conn.Close()
		return err
	}
	s.handlePacket(id, pkt)
	return nil
}
func (s *NIOServer) onPrepare(conn netpoll.Connection) context.Context {
	return context.Background()
}
func (s *NIOServer) onConnect(ctx context.Context, conn netpoll.Connection) context.Context {
	pkt, err := packet.ReadFrom(conn)
	if err != nil {
		conn.Close()
		return ctx
	}
	connPacket, ok := pkt.(*packet.ConnectPacket)
	if !ok {
		conn.Close()
		return ctx
	}
	s.conns.Write(func(m map[string]netpoll.Connection) {
		m[connPacket.Identity.ClientID] = conn
	})
	err = s.authHandler(connPacket.Identity)
	if err != nil {
		s.close(connPacket.Identity.ClientID, packet.CloseAuthenticationFailure)
		return ctx
	}
	s.sendPacket(connPacket.Identity.ClientID, packet.NewConnack(packet.ConnectionAccepted))
	return context.WithValue(ctx, identityKey, connPacket.Identity)
}
func (s *NIOServer) onDisconnect(ctx context.Context, conn netpoll.Connection) {
	id, ok := ctx.Value(identityKey).(*packet.Identity)
	if ok {
		s.conns.Delete(id.ClientID)
	}
}

type contextKey struct{}

var identityKey = contextKey{}
