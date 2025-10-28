package broker

import (
	"context"
	"strings"

	"sutext.github.io/entry/client"
)

type Peer interface {
	Join(channels []string, uid string) error
	Leave(channels []string, uid string) error
	SendData(ctx context.Context, channel string, payload []byte) error
}

type peer struct {
	client *client.Client
}

func NewPeer(client *client.Client) Peer {
	return &peer{
		client: client,
	}
}
func (p *peer) Join(channels []string, uid string) error {
	return p.client.SendData([]byte("join:" + uid + ":" + strings.Join(channels, ",")))
}
func (p *peer) Leave(channels []string, uid string) error {
	return p.client.SendData([]byte("leave:" + uid + ":" + strings.Join(channels, ",")))
}
func (p *peer) SendData(ctx context.Context, channel string, payload []byte) error {
	return p.client.SendData(payload)
}
