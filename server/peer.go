package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"sutext.github.io/entry/keepalive"
	"sutext.github.io/entry/packet"
)

type peer struct {
	mu         *sync.RWMutex
	raw        net.Conn
	ready      chan struct{}
	logger     *slog.Logger
	server     *Server
	keepaplive *keepalive.KeepAlive
}

func newPeer(raw net.Conn, server *Server) *peer {
	s := &peer{
		mu:         new(sync.RWMutex),
		raw:        raw,
		logger:     server.logger,
		server:     server,
		keepaplive: keepalive.New(server.config.KeepAlive, server.config.PingTimeout),
	}
	s.keepaplive.PingFunc(func() {
		s.SendPing()
	})
	s.keepaplive.TimeoutFunc(func() {
		s.Close(packet.CloseNormal)
	})
	return s
}
func (s *peer) ClientID() (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	addr := s.raw.RemoteAddr()
	if addr != nil {
		return addr.String(), true
	}
	return "", false
}

func (s *peer) Close(code packet.CloseCode) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.raw == nil {
		return
	}
	s.SendPacket(packet.Close(code))
	s.raw.Close()
	s.keepaplive.Stop()
	close(s.ready)
	s.mu = nil
	s.raw = nil
	s.keepaplive = nil
	s.server = nil
	s.logger = nil
	s.ready = nil
}
func (s *peer) isClosed() bool {
	return s.raw == nil
}
func (s *peer) serve() {
	go func() {
		timer := time.NewTimer(time.Second * 10)
		defer timer.Stop()
		select {
		case <-s.ready:
			return
		case <-timer.C:
			s.Close(packet.CloseAuthenticationTimeout)
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
				if s.isClosed() {
					return
				}
				packet, err := packet.ReadPacket(s.raw)
				if err != nil {
					s.logger.ErrorContext(ctx, "read packet failed", "error", err)
					return
				}
				go s.handlePacket(ctx, packet)
			}
		}
	}()
}
func (s *peer) Dataack(packetId int64) {
	s.SendPacket(packet.DataAck(packetId))
}
func (s *peer) Connack(code packet.ConnectCode) error {
	return s.SendPacket(packet.Connack(code))
}
func (s *peer) SendPing() error {
	return s.SendPacket(packet.Ping())
}
func (s *peer) SendPong() error {
	return s.SendPacket(packet.Pong())
}
func (s *peer) SendPacket(p packet.Packet) error {
	if s.raw == nil {
		return fmt.Errorf("connection already closed")
	}
	return packet.WritePacket(s.raw, p)
}

func (s *peer) handlePacket(ctx context.Context, p packet.Packet) {
	if s.isClosed() {
		return
	}
	s.logger.Debug("handle packet", "packet", p.String())
	switch p := p.(type) {
	case *packet.ConnectPacket:
		s.server.addPeer(s)
		s.Connack(packet.ConnectionAccepted)
	case *packet.DataPacket:
		s.server.handlePeerData(ctx, p)
	case *packet.PingPacket:
		s.SendPong()
	case *packet.PongPacket:
		s.keepaplive.HandlePong()
	default:
		break
	}
}
