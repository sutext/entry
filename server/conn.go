package server

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"sutext.github.io/entry/keepalive"
	"sutext.github.io/entry/packet"
)

type conn struct {
	mu         *sync.RWMutex
	raw        net.Conn
	user       *packet.Identity
	logger     *slog.Logger
	server     *Server
	loginOk    chan struct{}
	keepaplive *keepalive.KeepAlive
}

func newConn(raw net.Conn, server *Server) *conn {
	c := &conn{
		mu:         new(sync.RWMutex),
		raw:        raw,
		logger:     server.logger,
		server:     server,
		loginOk:    make(chan struct{}),
		keepaplive: keepalive.New(server.config.KeepAlive, server.config.PingTimeout),
	}
	c.keepaplive.PingFunc(func() {
		c.sendPing()
	})
	c.keepaplive.TimeoutFunc(func() {
		c.close(packet.CloseNormal)
	})
	return c
}
func (c *conn) clientId() (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.user != nil {
		return c.user.ClientID, true
	}
	if c.raw != nil {
		return c.raw.RemoteAddr().String(), true
	}
	return "", false
}

func (c *conn) close(code packet.CloseCode) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.raw == nil {
		return
	}
	c.logger.Info("Close connection", "code", code)
	c.sendPacket(packet.Close(code))
	c.raw.Close()
	c.keepaplive.Stop()
	close(c.loginOk)
	c.mu = nil
	c.raw = nil
	c.keepaplive = nil
	c.server = nil
	c.logger = nil
	c.user = nil
	c.loginOk = nil
}
func (c *conn) closed() bool {
	return c.raw == nil
}
func (c *conn) serve() {
	go func() {
		timer := time.NewTimer(time.Second * 10)
		defer timer.Stop()
		select {
		case <-c.loginOk:
			return
		case <-timer.C:
			c.close(packet.CloseAuthenticationTimeout)
			return
		}
	}()
	for {
		packet, err := packet.ReadPacket(c.raw)
		if err != nil {
			return
		}
		go c.handlePacket(packet)
	}
}

func (c *conn) connack(code packet.ConnectCode) error {
	return c.sendPacket(packet.Connack(code))
}
func (c *conn) sendPing() error {
	return c.sendPacket(packet.Ping())
}
func (c *conn) sendPong() error {
	return c.sendPacket(packet.Pong())
}
func (c *conn) sendPacket(p packet.Packet) error {
	if c.raw == nil {
		return fmt.Errorf("connection already closed")
	}
	return packet.WritePacket(c.raw, p)
}
func (c *conn) doLogin(p *packet.ConnectPacket) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed() {
		return fmt.Errorf("connection already closed")
	}
	if c.user != nil {
		return fmt.Errorf("already login")
	}
	err := c.server.loginHandler(p.Identity)
	if err != nil {
		return err
	}
	c.user = p.Identity
	c.loginOk <- struct{}{}
	go c.server.addConn(c)
	c.keepaplive.Start()
	c.logger.Info("Login success", "user_id", c.user.UserID, "client_id", c.user.ClientID)
	return nil
}
func (c *conn) handlePacket(p packet.Packet) {
	if c.closed() {
		return
	}
	c.logger.Debug("handle packet", "packet", p.String())
	switch p := p.(type) {
	case *packet.ConnectPacket:
		err := c.doLogin(p)
		if err != nil {
			c.logger.Error("login failed", "error", err)
			c.close(packet.CloseAuthenticationFailure)
			return
		}
		c.connack(packet.ConnectionAccepted)
	case *packet.DataPacket:
		res, err := c.server.dataHandler(p)
		if err != nil {
			c.logger.Error("data handler failed", "error", err)
			return
		}
		if res != nil {
			c.sendPacket(res)
		}
	case *packet.PingPacket:
		c.sendPong()
	case *packet.PongPacket:
		c.keepaplive.HandlePong()
	case *packet.ClosePacket:
		c.close(p.Code)
	default:
		break
	}
}
