package entry

import (
	"fmt"

	qc "golang.org/x/net/quic"
	"sutext.github.io/entry/broker"
	"sutext.github.io/entry/internal/server/bio"
	"sutext.github.io/entry/internal/server/nio"
	"sutext.github.io/entry/internal/server/quic"
	"sutext.github.io/entry/server"
)

func TCP() Transport {
	return &tcpTransport{"tcp"}
}
func NIOTCP() Transport {
	return &tcpTransport{"nio"}
}
func QUIC(config *qc.Config) Transport {
	return &quicTransport{config}
}

type Transport interface {
	Network() string
	quicConfig() *qc.Config
}

type tcpTransport struct {
	network string
}

func (t *tcpTransport) Network() string {
	return t.network
}
func (t *tcpTransport) quicConfig() *qc.Config {
	return nil
}

type quicTransport struct {
	config *qc.Config
}

func (t *quicTransport) Network() string {
	return "quic"
}
func (t *quicTransport) quicConfig() *qc.Config {
	return t.config
}
func NewServer(transport Transport, address string) (server.Server, error) {
	switch transport.Network() {
	case "tcp":
		return bio.NewBIO(address), nil
	case "nio":
		return nio.NewNIO(address), nil
	case "quic":
		return quic.NewQUIC(address, transport.quicConfig()), nil
	default:
		return nil, fmt.Errorf("")
	}
}
func NewBroker(config *broker.Config) broker.Broker {
	return broker.New()
}
