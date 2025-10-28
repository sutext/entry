package broker

import (
	"context"

	"sutext.github.io/entry/internal/queue"
	"sutext.github.io/entry/internal/safe"
	"sutext.github.io/entry/internal/server/bio"
	"sutext.github.io/entry/internal/server/nio"
	"sutext.github.io/entry/packet"
	"sutext.github.io/entry/server"
)

type Broker interface {
	Join(channels []string, uid string) error
	Leave(channels []string, uid string) error
	SendData(ctx context.Context, channel string, payload []byte) error
	Start() error
	Shutdown() error
}

type broker struct {
	peers      *safe.Map[string, Peer]
	config     *Config
	channels   *safe.Map[string, *safe.Set[string]]
	taskQueue  *queue.Queue
	userServer server.Server
	peerServer server.Server
}

func New() Broker {
	conf := DefaultConfig()
	return &broker{
		config:    conf,
		peers:     safe.NewMap[string, Peer](),
		taskQueue: queue.NewQueue(10, 20),
		channels:  safe.NewMap[string, *safe.Set[string]](),
	}
}

func (b *broker) Start() error {
	b.userServer = bio.NewBIO(b.config.Addr)
	b.peerServer = nio.NewNIO(b.config.Peer.Addr)
	b.userServer.OnAuth(func(p *packet.Identity) error {
		return nil
	})
	b.userServer.OnData(func(uid string, p *packet.DataPacket) (*packet.DataPacket, error) {
		return nil, nil
	})
	return b.userServer.Serve()
}

func (b *broker) Shutdown() error {
	return b.userServer.Shutdown(context.Background())
}

func (b *broker) SendData(ctx context.Context, channel string, payload []byte) error {
	if users, ok := b.channels.Get(channel); ok {
		users.Range(func(uid string) bool {
			b.taskQueue.Push(func() {
				if conn, err := b.userServer.GetConn(uid); err != nil {
					conn.SendData(payload)
				}
			})
			return true
		})
	}
	b.peers.Range(func(key string, peer Peer) bool {
		peer.SendData(ctx, channel, payload)
		return true
	})
	return nil
}

func (b *broker) Join(channels []string, uid string) error {
	for _, channel := range channels {
		if users, ok := b.channels.Get(channel); ok {
			users.Add(uid)
		} else {
			b.channels.Set(channel, safe.NewSet(uid))
		}
	}
	b.peers.Range(func(key string, peer Peer) bool {
		peer.Join(channels, uid)
		return true
	})
	return nil
}

func (b *broker) Leave(channels []string, uid string) error {
	for _, channel := range channels {
		if users, ok := b.channels.Get(channel); ok {
			users.Del(uid)
		}
	}
	b.peers.Range(func(key string, peer Peer) bool {
		peer.Leave(channels, uid)
		return true
	})
	return nil
}
