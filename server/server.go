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

	"sutext.github.io/entry/logger"
	"sutext.github.io/entry/packet"
)

type DataHandler func(ctx context.Context, conn *Conn, packet packet.DataPacket) error
type LoginHandler func(ctx context.Context, conn *Conn, packet packet.ConnectPacket) error
type KickInfo struct {
	clientOK   bool
	oldClients map[string]*Conn
	newClient  *Conn
}
type Server struct {
	mu           *sync.Mutex
	conns        map[string]*Conn
	addChan      chan *Conn
	delChan      chan *Conn
	kickChan     chan *KickInfo
	config       *Config
	logger       *slog.Logger
	dataHandler  DataHandler
	loginHandler LoginHandler
}

func New(config *Config) *Server {
	s := &Server{
		mu:       &sync.Mutex{},
		conns:    make(map[string]*Conn),
		addChan:  make(chan *Conn),
		delChan:  make(chan *Conn),
		kickChan: make(chan *KickInfo),
		config:   config,
		logger:   logger.New(config.LoggerLevel, config.LoggerFormat),
		loginHandler: func(ctx context.Context, conn *Conn, packet packet.ConnectPacket) error {
			return fmt.Errorf("login handler not set")
		},
		dataHandler: func(ctx context.Context, conn *Conn, packet packet.DataPacket) error {
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
		for {
			select {
			case <-ctx.Done():
				return
			case conn := <-s.addChan:
				s.register(conn)
			case conn := <-s.delChan:
				s.unregister(conn)
			case kickout := <-s.kickChan:
				s.kickout(kickout.clientOK, kickout.oldClients, kickout.newClient)
			}
		}
	}()
	go func() {
		defer close(done)
		<-ctx.Done()
		s.Shutdown(context.Background())
	}()

	<-ctx.Done()
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
	s.conns[conn.cid] = conn
	s.mu.Unlock()
}
func (s *Server) unregister(conn *Conn) {
	s.mu.Lock()
	delete(s.conns, conn.cid)
	s.mu.Unlock()
}
func (s *Server) kickout(clientOK bool, oldClients map[string]*Conn, newClient *Conn) {

}
