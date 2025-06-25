package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"sutext.github.io/entry/keepalive"
	"sutext.github.io/entry/packet"
)

type Conn struct {
	raw        net.Conn
	cid        string
	server     *Server
	keepaplive *keepalive.KeepAlive
}

func newConn(raw net.Conn, server *Server) *Conn {
	c := &Conn{
		raw:        raw,
		server:     server,
		keepaplive: keepalive.New(server.config.KeepAlive, server.config.PingTimeout),
	}
	c.keepaplive.PingFunc(func() {
		c.SendPacket(packet.Ping())
	})
	c.keepaplive.TimeoutFunc(func() {
		c.server.delChan <- c
	})
	return c
}
func (c *Conn) start() {
	ctx := context.Background()
	timer := time.NewTimer(time.Second * 10)
	go func() {
		defer c.raw.Close()
		defer timer.Stop()
		select {
		case <-timer.C:
			c.server.logger.WarnContext(ctx, "login timeout")
			return
		case <-ctx.Done():
			c.server.logger.DebugContext(ctx, "connection closed", "error", ctx.Err())
			return
		}
	}()
	go func() {
		ctx, cancel := context.WithCancelCause(ctx)
		for {
			packet, err := packet.ReadPacket(c.raw)
			if err != nil {
				cancel(err)
				c.server.logger.ErrorContext(ctx, "read packet failed", "error", err)
				return
			}
			c.handlePacket(ctx, packet, timer)
		}
	}()
}
func (c *Conn) SendPacket(p packet.Packet) error {
	return packet.WritePacket(c.raw, p)
}
func (c *Conn) handlePacket(ctx context.Context, p packet.Packet, timer *time.Timer) {
	c.server.logger.DebugContext(ctx, "handle packet", "packet", p.String())
	switch p := p.(type) {
	case *packet.ConnectPacket:
		c.cid = fmt.Sprintf("%s/%d", p.UserID, p.Platform)
		err := c.server.loginHandler(ctx, c, *p)
		if err != nil {
			c.server.logger.ErrorContext(ctx, "login failed", "error", err)
			return
		}
		timer.Stop()
		c.server.addChan <- c
		c.SendPacket(packet.Connack(0))
		c.keepaplive.Start()
		c.server.logger.InfoContext(ctx, "new connection", "user_id", p.UserID, "platform", p.Platform)
	case *packet.DataPacket:
		err := c.server.dataHandler(ctx, c.server.conns[c.cid], *p)
		if err != nil {
			c.server.logger.ErrorContext(ctx, "data handler failed", "error", err)
			return
		}
		if p.Qos > 0 {
			c.SendPacket(packet.DataAck(p.PacketId))
		}
	case *packet.PingPacket:
		c.SendPacket(packet.Pong())
	case *packet.PongPacket:
		c.keepaplive.HandlePong()
	default:
		break
	}
}
