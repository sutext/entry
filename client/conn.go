package client

import (
	"net"

	"sutext.github.io/entry/packet"
)

type conn struct {
	raw     net.Conn
	pkgFunc func(packet.Packet)
	errFunc func(error)
}

func (c *conn) onError(f func(error)) {
	c.errFunc = f
}
func (c *conn) onPacket(f func(packet.Packet)) {
	c.pkgFunc = f
}
func (c *conn) connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	c.raw = conn
	go c.receive()
	return nil
}
func (c *conn) receive() {
	for {
		p, err := packet.ReadPacket(c.raw)
		if err != nil {
			if c.errFunc != nil {
				c.errFunc(err)
			}
			return
		}
		if c.pkgFunc != nil {
			c.pkgFunc(p)
		}
	}
}
func (c *conn) close() error {
	if c.raw == nil {
		return nil
	}
	c.pkgFunc = nil
	c.errFunc = nil
	return c.raw.Close()
}

func (c *conn) sendPacket(p packet.Packet) error {
	if c.raw == nil {
		return ErrNotConnected
	}
	return packet.WritePacket(c.raw, p)
}
