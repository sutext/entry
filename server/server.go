package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"sutext.github.io/entry/code"
	"sutext.github.io/entry/logger"
	"sutext.github.io/entry/packet"
)

type DataHandler func(ctx context.Context, conn *Conn, p *packet.DataPacket) error
type LoginHandler func(ctx context.Context, conn *Conn, p *packet.ConnectPacket) error
type ConnMap map[code.Platform]*Conn

type Server struct {
	mu           *sync.RWMutex
	conns        map[string]ConnMap
	config       *Config
	logger       *slog.Logger
	dataHandler  DataHandler
	loginHandler LoginHandler
}

func New(config *Config) *Server {
	s := &Server{
		mu:     &sync.RWMutex{},
		conns:  make(map[string]ConnMap),
		config: config,
		logger: logger.New(config.LoggerLevel, config.LoggerFormat),
		loginHandler: func(ctx context.Context, conn *Conn, p *packet.ConnectPacket) error {
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
				newConn(conn, s).start()
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
func (s *Server) SendPacket(ctx context.Context, uid string, packet packet.Packet) {

}
func (s *Server) register(conn *Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user := conn.User()
	if user == nil {
		return
	}
	if mp, ok := s.conns[user.UserID]; ok {
		if old, ok := mp[user.Platform]; ok {
			old.Close(packet.CloseDuplicateLogin)
		}
		mp[user.Platform] = conn
	} else {
		s.conns[user.UserID] = ConnMap{user.Platform: conn}
	}
}

func (s *Server) KickUser(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if mp, ok := s.conns[userID]; ok {
		for _, conn := range mp {
			conn.Close(packet.CloseKickedOut)
		}
		delete(s.conns, userID)
	}
}
func (s *Server) Kickout(clientOK bool, oldClients map[string]*Conn, newClient *Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
}
func (s *Server) GetConn(userID string, platform code.Platform) *Conn {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if mp, ok := s.conns[userID]; ok {
		if conn, ok := mp[platform]; ok {
			return conn
		}
	}
	return nil
}
func (s *Server) GetConns(userID string) ConnMap {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if mp, ok := s.conns[userID]; ok {
		return mp
	}
	return nil
}
