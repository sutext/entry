package server

import (
	"context"
	"errors"

	"sutext.github.io/entry/packet"
)

type Conn interface {
	GetID() *packet.Identity
	Close(reason packet.CloseCode)
	SendPing() error
	SendPong() error
	SendData(data []byte) error
}
type OnData func(cid string, p *packet.DataPacket) (*packet.DataPacket, error)
type OnAuth func(p *packet.Identity) error

type Server interface {
	Serve() error
	OnAuth(handler OnAuth)
	OnData(handler OnData)
	GetConn(cid string) (Conn, error)
	KickConn(cid string) error
	Shutdown(ctx context.Context) error
}

var (
	ErrConnClosed   = errors.New("connection_closed")
	ErrAuthFailed   = errors.New("auth_failed")
	ErrConnNotFound = errors.New("connection_not_found")
)
