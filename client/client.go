package client

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"sync"

	"sutext.github.io/entry/keepalive"
	"sutext.github.io/entry/logger"
	"sutext.github.io/entry/packet"
	"sutext.github.io/entry/types"
)

var ErrNotConnected = errors.New("not connected")

type Identity struct {
	UserID      string
	AccessToken string
}
type Client struct {
	mu        *sync.Mutex
	conn      net.Conn
	host      string
	prot      string
	status    Status
	logger    *slog.Logger
	identity  *Identity
	platform  types.Platform
	keepalive *keepalive.KeepAlive
}

func New(config *Config) *Client {
	c := &Client{
		mu:        new(sync.Mutex),
		host:      config.Host,
		prot:      config.Port,
		status:    StatusUnknown,
		logger:    logger.New(config.LoggerLevel, config.LoggerFormat),
		platform:  config.Platform,
		keepalive: keepalive.New(config.KeepAlive, config.PingTimeout),
	}
	c.keepalive.PingFunc(func() {
		c.SendPacket(packet.Ping())
	})
	c.keepalive.TimeoutFunc(func() {
		c.logger.Error("keepalive timeout")
		go c.reconnect()
	})
	return c
}
func (c *Client) Connect(userId string, accessToken string) {
	c.identity = &Identity{
		UserID:      userId,
		AccessToken: accessToken,
	}
	go c.reconnect()
}
func (c *Client) reconnect() {
	if c.identity == nil {
		return
	}
	if c.Status() == StatusOpened || c.Status() == StatusOpening {
		return
	}
	conn, err := net.Dial("tcp", net.JoinHostPort(c.host, c.prot))
	if err != nil {
		c.logger.Error(err.Error())
		c.setStatus(StatusClosed)
		return
	}
	c.conn = conn
	c.setStatus(StatusOpening)
	err = c.SendPacket(packet.Connect(c.identity.UserID, c.platform, c.identity.AccessToken))
	if err != nil {
		c.logger.Error(err.Error())
		c.setStatus(StatusClosed)
		return
	}
	for {
		p, err := packet.ReadPacket(conn)
		if err != nil {
			c.logger.Error(err.Error())
			c.setStatus(StatusClosed)
			return
		}
		c.handlePacket(p)
	}
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
	return packet.WritePacket(c.conn, p)
}
func (c *Client) Status() Status {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.status
}
func (c *Client) setStatus(status Status) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.status == status {
		return
	}
	c.status = status
	if status == StatusClosed {
		if c.conn != nil {
			c.conn.Close()
			c.conn = nil
		}
		c.keepalive.Stop()
		go c.reconnect()
	}
	// c.NotifyStatus <- status
}
func (c *Client) handlePacket(p packet.Packet) {
	c.logger.Info("receive packet", "packet", p.String())
	switch p := p.(type) {
	case *packet.ConnackPacket:
		if p.Code != 0 {
			c.logger.Error("connect failed")
		}
		c.setStatus(StatusOpened)
		c.keepalive.Start()
	case *packet.DataPacket:
	case *packet.PingPacket:
		c.SendPacket(packet.Pong())
	case *packet.PongPacket:
		c.keepalive.HandlePong()
	case *packet.ClosePacket:
		c.setStatus(StatusClosed)
	default:

	}
}
