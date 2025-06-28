package client

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"sync"
	"time"

	"sutext.github.io/entry/backoff"
	"sutext.github.io/entry/keepalive"
	"sutext.github.io/entry/logger"
	"sutext.github.io/entry/packet"
)

var ErrNotConnected = errors.New("not connected")

type Client struct {
	mu        *sync.RWMutex
	conn      *conn
	config    *Config
	status    Status
	logger    *slog.Logger
	retrier   *Retrier
	identity  *packet.Identity
	retrying  bool
	keepalive *keepalive.KeepAlive
}

func New(config *Config) *Client {
	c := &Client{
		mu:        new(sync.RWMutex),
		config:    config,
		status:    StatusUnknown,
		logger:    logger.New(config.LoggerLevel, config.LoggerFormat),
		retrier:   NewRetrier(100000, backoff.Constant(time.Second*2)),
		keepalive: keepalive.New(config.KeepAlive, config.PingTimeout),
	}
	c.keepalive.PingFunc(func() {
		c.SendPacket(packet.Ping())
	})
	c.keepalive.TimeoutFunc(func() {
		c.logger.Error("keepalive timeout")
		c.tryClose(CloseReasonPingTimeout)
	})
	return c
}

func (c *Client) Connect(identity *packet.Identity) {
	c.identity = identity
	switch c.Status() {
	case StatusOpened, StatusOpening:
		return
	}
	c.setStatus(StatusOpening)
	c.reconnect()
}
func (c *Client) tryClose(err error) {
	c.logger.Error("try close", "reason", err)
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.identity == nil {
		return
	}
	if c.status == StatusClosed || c.status == StatusClosing {
		return
	}
	if c.retrying {
		return
	}
	if code, ok := err.(CloseReason); ok {
		if code == CloseReasonNormal {
			c.setStatus(StatusClosed)
			return
		}
	}
	if c.retrier == nil {
		c.setStatus(StatusClosed)
		return
	}
	delay, ok := c.retrier.can(err)
	if !ok {
		c.setStatus(StatusClosed)
		return
	}
	c.retrying = true
	c.setStatus(StatusOpening)
	c.logger.Info("will retry after", "delay", delay.String())
	c.retrier.retry(delay, func() {
		c.retrying = false
		c.reconnect()
	})

}
func (c *Client) reconnect() {
	if c.conn != nil {
		c.conn.close()
		c.conn = nil
	}
	c.conn = &conn{}
	c.conn.onPacket(func(p packet.Packet) {
		c.handlePacket(context.Background(), p)
	})
	c.conn.onError(func(err error) {
		c.tryClose(err)
	})
	err := c.conn.connect(net.JoinHostPort(c.config.Host, c.config.Port))
	if err != nil {
		c.tryClose(err)
		return
	}
	c.SendPacket(packet.Connect(c.identity))
}

func (c *Client) SendData(data []byte, packetId int64, dataType packet.DataType) error {
	dataPacket := packet.Data(dataType, packetId, data)
	return c.SendPacket(dataPacket)
}
func (c *Client) SendData0(data []byte, dataType packet.DataType) error {
	dataPacket := packet.Data0(dataType, data)
	return c.SendPacket(dataPacket)
}
func (c *Client) SendText(text string, packetId int64) error {
	dataPacket := packet.Data(packet.DataTypeText, packetId, []byte(text))
	return c.SendPacket(dataPacket)
}
func (c *Client) SendText0(text string) error {
	dataPacket := packet.Data0(packet.DataTypeText, []byte(text))
	return c.SendPacket(dataPacket)
}
func (c *Client) SendJSON(j any, packetId int64) error {
	jsonData, err := json.Marshal(j)
	if err != nil {
		return err
	}
	dataPacket := packet.Data(packet.DataTypeText, packetId, jsonData)
	return c.SendPacket(dataPacket)
}
func (c *Client) SendPing() error {
	return c.SendPacket(packet.Ping())
}
func (c *Client) SendPong() error {
	return c.SendPacket(packet.Pong())
}
func (c *Client) SendJSON0(j any) error {
	jsonData, err := json.Marshal(j)
	if err != nil {
		return err
	}
	dataPacket := packet.Data0(packet.DataTypeText, jsonData)
	return c.SendPacket(dataPacket)
}
func (c *Client) SendPacket(p packet.Packet) error {
	if c.conn == nil {
		return ErrNotConnected
	}
	return c.conn.sendPacket(p)
}
func (c *Client) Status() Status {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status
}
func (c *Client) safeSetStatus(status Status) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.setStatus(status)
}
func (c *Client) setStatus(status Status) {
	if c.status == status {
		return
	}
	c.logger.Info("status change", "from", c.status.String(), "to", status.String())
	c.status = status
	switch status {
	case StatusClosed:
		c.keepalive.Stop()
		c.retrier.cancel()
		if c.conn != nil {
			c.conn.close()
			c.conn = nil
		}
	case StatusOpening, StatusClosing:
		c.keepalive.Stop()
	case StatusOpened:
		c.keepalive.Start()
	}
}
func (c *Client) handlePacket(ctx context.Context, p packet.Packet) {
	c.logger.Info("receive packet", "packet", p.String())
	switch p := p.(type) {
	case *packet.ConnackPacket:
		if p.Code != 0 {
			return
		}
		c.safeSetStatus(StatusOpened)
	case *packet.DataPacket:

	case *packet.PingPacket:
		c.SendPong()
	case *packet.PongPacket:
		c.keepalive.HandlePong()
	case *packet.ClosePacket:
		c.tryClose(p.Code)
	default:

	}
}
