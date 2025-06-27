package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"sutext.github.io/entry/code"
	"sutext.github.io/entry/keepalive"
	"sutext.github.io/entry/packet"
)

type ConnUser struct {
	UserID      string
	AccessToken string
	Platform    code.Platform
}
type Conn struct {
	mu         *sync.RWMutex
	raw        net.Conn
	user       *ConnUser
	logger     *slog.Logger
	server     *Server
	loginOk    chan struct{}
	keepaplive *keepalive.KeepAlive
}

func newConn(raw net.Conn, server *Server) *Conn {
	c := &Conn{
		mu:         new(sync.RWMutex),
		raw:        raw,
		logger:     server.logger,
		server:     server,
		loginOk:    make(chan struct{}),
		keepaplive: keepalive.New(server.config.KeepAlive, server.config.PingTimeout),
	}
	c.keepaplive.PingFunc(func() {
		c.SendPing()
	})
	c.keepaplive.TimeoutFunc(func() {
		c.Close(packet.CloseNormal)
	})
	return c
}
func (c *Conn) User() *ConnUser {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.user
}
func (c *Conn) Close(code packet.CloseCode) {
	if c.raw == nil {
		return
	}
	c.SendPacket(packet.Close(code))
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
func (c *Conn) isClosed() bool {
	return c.raw == nil
}
func (c *Conn) start() {
	go func() {
		timer := time.NewTimer(time.Second * 10)
		defer timer.Stop()
		select {
		case <-c.loginOk:
			return
		case <-timer.C:
			c.Close(packet.CloseAuthenticationTimeout)
			return
		}
	}()
	go func() {
		ctx := context.Background()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if c.isClosed() {
					return
				}
				packet, err := packet.ReadPacket(c.raw)
				if err != nil {
					c.logger.ErrorContext(ctx, "read packet failed", "error", err)
					return
				}
				go c.handlePacket(ctx, packet)
			}
		}
	}()
}
func (c *Conn) Dataack(packetId int64) {
	c.SendPacket(packet.DataAck(packetId))
}
func (c *Conn) Connack(code packet.ConnectCode) error {
	return c.SendPacket(packet.Connack(code))
}
func (c *Conn) SendPing() error {
	return c.SendPacket(packet.Ping())
}
func (c *Conn) SendPong() error {
	return c.SendPacket(packet.Pong())
}
func (c *Conn) SendPacket(p packet.Packet) error {
	if c.raw == nil {
		return fmt.Errorf("connection already closed")
	}
	return packet.WritePacket(c.raw, p)
}
func (c *Conn) doLogin(ctx context.Context, p *packet.ConnectPacket) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isClosed() {
		return fmt.Errorf("connection already closed")
	}
	if c.user != nil {
		return fmt.Errorf("already login")
	}
	err := c.server.loginHandler(ctx, c, p)
	if err != nil {
		return err
	}
	c.user = &ConnUser{
		UserID:      p.UserID,
		AccessToken: p.AccessToken,
		Platform:    p.Platform,
	}
	c.loginOk <- struct{}{}
	go c.server.register(c)
	c.keepaplive.Start()
	c.logger.Info("Login success", "user_id", p.UserID, "platform", p.Platform)
	return nil
}
func (c *Conn) handlePacket(ctx context.Context, p packet.Packet) {
	if c.isClosed() {
		return
	}
	c.logger.Debug("handle packet", "packet", p.String())
	switch p := p.(type) {
	case *packet.ConnectPacket:
		err := c.doLogin(ctx, p)
		if err != nil {
			c.logger.Error("login failed", "error", err)
			cancelContext(ctx, err)
			return
		}
		c.Connack(packet.ConnectionAccepted)
	case *packet.DataPacket:
		err := c.server.dataHandler(ctx, c, p)
		if err != nil {
			c.logger.Error("data handler failed", "error", err)
			return
		}
		if p.Qos > 0 {
			c.Dataack(p.PacketId)
		}
	case *packet.PingPacket:
		c.SendPong()
	case *packet.PongPacket:
		c.keepaplive.HandlePong()
	default:
		break
	}
}
func cancelContext(ctx context.Context, err error) {
	ctx, cancel := context.WithCancelCause(ctx)
	cancel(err)
}
