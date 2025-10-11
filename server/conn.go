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
	mu        *sync.RWMutex
	raw       net.Conn
	auth      *packet.Identity
	logger    *slog.Logger
	server    *Server
	authed    chan struct{}
	keepAlive *keepalive.KeepAlive
}

func newConn(raw net.Conn, server *Server) *conn {
	c := &conn{
		mu:     new(sync.RWMutex),
		raw:    raw,
		logger: server.logger,
		server: server,
		authed: make(chan struct{}),
	}
	if server.keepAlive != nil {
		c.keepAlive = keepalive.New(server.keepAlive.interval, server.keepAlive.timeout)
		c.keepAlive.PingFunc(func() {
			c.sendPing()
		})
		c.keepAlive.TimeoutFunc(func() {
			c.close(packet.CloseNormal)
		})
	}
	return c
}
func (c *conn) clientId() (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.auth != nil {
		return c.auth.ClientID, true
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
	c.sendPacket(packet.NewClose(code))
	c.raw.Close()
	c.keepAlive.Stop()
	close(c.authed)
	c.mu = nil
	c.raw = nil
	c.keepAlive = nil
	c.server = nil
	c.logger = nil
	c.auth = nil
	c.authed = nil
}
func (c *conn) closed() bool {
	return c.raw == nil
}
func (c *conn) serve() {
	go func() {
		timer := time.NewTimer(time.Second * 10)
		defer timer.Stop()
		select {
		case <-c.authed:
			return
		case <-timer.C:
			c.close(packet.CloseAuthenticationTimeout)
			return
		}
	}()
	for {
		packet, err := packet.ReadFrom(c.raw)
		if err != nil {
			return
		}
		go c.handlePacket(packet)
	}
}

func (c *conn) connack(code packet.ConnectCode) error {
	return c.sendPacket(packet.NewConnack(code))
}
func (c *conn) sendPing() error {
	return c.sendPacket(packet.NewPing())
}
func (c *conn) sendPong() error {
	return c.sendPacket(packet.NewPong())
}
func (c *conn) sendPacket(p packet.Packet) error {
	if c.raw == nil {
		return fmt.Errorf("connection already closed")
	}
	return packet.WriteTo(c.raw, p)
}
func (c *conn) doAuth(id *packet.Identity) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed() {
		return fmt.Errorf("connection already closed")
	}
	if c.auth != nil {
		return fmt.Errorf("already login")
	}
	err := c.server.authHandler(id)
	if err != nil {
		return err
	}
	c.auth = id
	c.authed <- struct{}{}
	go c.server.addConn(c)
	c.keepAlive.Start()
	c.logger.Info("Login success", "user_id", c.auth.UserID, "client_id", c.auth.ClientID)
	return nil
}
func (c *conn) handlePacket(p packet.Packet) {
	if c.closed() {
		return
	}
	// c.logger.Debug("handle packet", "packet", p.String())
	switch p.Type() {
	case packet.CONNECT:
		p := p.(*packet.ConnectPacket)
		if p.Identity != nil {
			err := c.doAuth(p.Identity)
			if err != nil {
				c.logger.Error("login failed", "error", err)
				c.close(packet.CloseAuthenticationFailure)
				return
			}
		}
		c.connack(packet.ConnectionAccepted)
	case packet.DATA:
		p := p.(*packet.DataPacket)
		res, err := c.server.dataHandler(p)
		if err != nil {
			c.logger.Error("data handler failed", "error", err)
			return
		}
		if res != nil {
			c.sendPacket(res)
		}
	case packet.PING:
		c.sendPong()
	case packet.PONG:
		c.keepAlive.HandlePong()
	case packet.CLOSE:
		c.close(p.(*packet.ClosePacket).Code)
	default:
		break
	}
}
