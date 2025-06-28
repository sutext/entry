package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sutext.github.io/entry/logger"
	"sutext.github.io/entry/packet"
	"sutext.github.io/entry/safe"
)

type DataHandler func(ctx context.Context, conn *Conn, p *packet.DataPacket) error
type LoginHandler func(ctx context.Context, conn *Conn, p *packet.Identity) error
type Server struct {
	subs         *safe.Map[string, []*Conn]
	conns        *safe.Map[string, *Conn]
	peers        *safe.Map[string, *peer]
	config       *Config
	logger       *slog.Logger
	dataHandler  DataHandler
	loginHandler LoginHandler
}

func New(config *Config) *Server {
	s := &Server{
		subs:   safe.NewMap(map[string][]*Conn{}),
		conns:  safe.NewMap(map[string]*Conn{}),
		peers:  safe.NewMap(map[string]*peer{}),
		config: config,
		logger: logger.New(config.LoggerLevel, config.LoggerFormat),
		loginHandler: func(ctx context.Context, conn *Conn, p *packet.Identity) error {
			return fmt.Errorf("login handler not set")
		},
		dataHandler: func(ctx context.Context, conn *Conn, p *packet.DataPacket) error {
			return fmt.Errorf("data handler not set")
		},
	}
	return s
}

func (s *Server) Run(ctx context.Context) error {
	s.logger.InfoContext(ctx, "entry server start")
	ctx, cancel := context.WithCancelCause(ctx)
	done := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel(fmt.Errorf("entry signal received"))
	}()
	go func() {
		listener, err := net.Listen("tcp", ":"+s.config.Port)
		if err != nil {
			cancel(fmt.Errorf("entry %w", err))
			return
		}
		for {
			select {
			case <-ctx.Done():
				fmt.Println("entry server stoped1")
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					cancel(fmt.Errorf("entry %w", err))
					return
				}
				c := newConn(conn, s)
				c.serve()
			}
		}
	}()
	go func() {
		listener, err := net.Listen("tcp", ":4567")
		if err != nil {
			cancel(fmt.Errorf("entry %w", err))
			return
		}
		for {
			select {
			case <-ctx.Done():
				fmt.Println("entry server stoped1")
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					cancel(fmt.Errorf("entry %w", err))
					return
				}
				c := newPeer(conn, s)
				c.serve()
			}
		}
	}()
	go func() {
		defer close(done)
		<-ctx.Done()
		fmt.Println("entry server stoped2")
		s.Shutdown(context.Background())
	}()

	<-ctx.Done()
	fmt.Println("entry server stoped3")
	s.logger.InfoContext(ctx, "entry server stoped")
	timeout := time.NewTimer(time.Second * 15)
	defer timeout.Stop()
	select {
	case <-timeout.C:
		s.logger.WarnContext(ctx, "entry server graceful shutdown timeout")
	case <-done:
		s.logger.DebugContext(ctx, "entry server graceful shutdown")
	}
	return context.Cause(ctx)
}
func (s *Server) HandleLogin(handler LoginHandler) {
	s.loginHandler = handler
}
func (s *Server) HandleData(handler DataHandler) {
	s.dataHandler = handler
}
func (s *Server) Shutdown(ctx context.Context) {
}
func (s *Server) SendData(p *packet.DataPacket, to string) error {
	if conn, ok := s.conns.Get(to); ok {
		return conn.SendPacket(p)
	}
	return fmt.Errorf("conn not found")
}
func (s *Server) handlePeerData(ctx context.Context, p *packet.DataPacket) {
	if conns, ok := s.subs.Get(""); ok {
		for _, conn := range conns {
			conn.SendPacket(p)
		}
	}
}
func (s *Server) register(conn *Conn) {
	clientId, ok := conn.ClientID()
	if !ok {
		return
	}
	s.conns.Write(func(m map[string]*Conn) {
		if old, ok := m[clientId]; ok {
			old.Close(packet.CloseDuplicateLogin)
		}
		m[clientId] = conn
	})
}
func (s *Server) addPeer(peer *peer) {
	clientId, ok := peer.ClientID()
	if !ok {
		return
	}
	s.peers.Set(clientId, peer)
}

func (s *Server) KickClient(clientID string) {
	s.conns.Write(func(m map[string]*Conn) {
		if conn, ok := m[clientID]; ok {
			conn.Close(packet.CloseKickedOut)
			delete(m, clientID)
		}
	})
}
