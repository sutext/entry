package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"sutext.github.io/entry/logger"
	"sutext.github.io/entry/packet"
	"sutext.github.io/entry/safe"
)

type DataHandler func(p *packet.DataPacket) (*packet.DataPacket, error)
type AuthHandler func(p *packet.Identity) error
type Server struct {
	conns     *safe.Map[string, *conn]
	logger    *slog.Logger
	keepAlive *struct {
		interval time.Duration
		timeout  time.Duration
	}
	dataHandler DataHandler
	authHandler AuthHandler
}

func New() *Server {
	s := &Server{
		conns:  safe.NewMap(map[string]*conn{}),
		logger: logger.New(logger.LevelDebug, logger.FormatJSON),
		authHandler: func(p *packet.Identity) error {
			return fmt.Errorf("login handler not set")
		},
		dataHandler: func(p *packet.DataPacket) (*packet.DataPacket, error) {
			return nil, fmt.Errorf("data handler not set")
		},
	}
	return s
}
func (s *Server) SetLogger(logger *slog.Logger) {
	s.logger = logger
}
func (s *Server) SetKeepAlive(interval time.Duration, timeout time.Duration) {
	s.keepAlive = &struct {
		interval time.Duration
		timeout  time.Duration
	}{
		interval: interval,
		timeout:  timeout,
	}
}

func (s *Server) Listen(ctx context.Context, port string) {
	ctx, cancel := context.WithCancelCause(ctx)
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		cancel(err)
		return
	}
	go func() {
		<-ctx.Done()
		listener.Close()
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			cancel(err)
			return
		}
		c := newConn(conn, s)
		go c.serve()
	}
}
func (s *Server) HandleAuth(handler AuthHandler) {
	s.authHandler = handler
}
func (s *Server) HandleData(handler DataHandler) {
	s.dataHandler = handler
}
func (s *Server) Shutdown(ctx context.Context) {
}
func (s *Server) SendData(data []byte, clientID string) error {
	if conn, ok := s.conns.Get(clientID); ok {
		return conn.sendPacket(packet.NewData(data))
	}
	return fmt.Errorf("conn not found")
}
func (s *Server) KickClient(clientID string) {
	s.conns.Write(func(m map[string]*conn) {
		if conn, ok := m[clientID]; ok {
			conn.close(packet.CloseKickedOut)
			delete(m, clientID)
		}
	})
}
func (s *Server) addConn(c *conn) {
	clientId, ok := c.clientId()
	if !ok {
		return
	}
	s.conns.Write(func(m map[string]*conn) {
		if old, ok := m[clientId]; ok {
			old.close(packet.CloseDuplicateLogin)
		}
		m[clientId] = c
	})
}
