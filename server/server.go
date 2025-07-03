package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"sutext.github.io/entry/logger"
	"sutext.github.io/entry/packet"
	"sutext.github.io/entry/safe"
)

type DataHandler func(p *packet.DataPacket) (*packet.DataPacket, error)
type LoginHandler func(p *packet.Identity) error
type Server struct {
	conns        *safe.Map[string, *conn]
	config       *Config
	logger       *slog.Logger
	dataHandler  DataHandler
	loginHandler LoginHandler
}

func New(config *Config) *Server {
	s := &Server{
		conns:  safe.NewMap(map[string]*conn{}),
		config: config,
		logger: logger.New(config.LoggerLevel, config.LoggerFormat),
		loginHandler: func(p *packet.Identity) error {
			return fmt.Errorf("login handler not set")
		},
		dataHandler: func(p *packet.DataPacket) (*packet.DataPacket, error) {
			return nil, fmt.Errorf("data handler not set")
		},
	}
	return s
}
func (s *Server) Listen(ctx context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)
	listener, err := net.Listen("tcp", ":"+s.config.Port)
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
func (s *Server) HandleLogin(handler LoginHandler) {
	s.loginHandler = handler
}
func (s *Server) HandleData(handler DataHandler) {
	s.dataHandler = handler
}
func (s *Server) Shutdown(ctx context.Context) {
}
func (s *Server) SendData(data []byte, clientID string) error {
	if conn, ok := s.conns.Get(clientID); ok {
		return conn.sendPacket(packet.Data(packet.DataBinary, data))
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
