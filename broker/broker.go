package broker

import (
	"context"

	"sutext.github.io/entry/internal/queue"
	"sutext.github.io/entry/internal/safe"
	"sutext.github.io/entry/internal/server/bio"
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
	config    *Config
	server    server.Server
	channels  *safe.Map[string, *safe.Set[string]]
	taskQueue *queue.Queue
}

func New() Broker {
	conf := DefaultConfig()
	return &broker{
		config:    conf,
		taskQueue: queue.NewQueue(10, 20),
		channels:  safe.NewMap[string, *safe.Set[string]](),
	}
}

func (b *broker) Start() error {
	b.server = bio.NewBIO(b.config.Addr)
	b.server.OnAuth(func(p *packet.Identity) error {
		return nil
	})
	b.server.OnData(func(uid string, p *packet.DataPacket) (*packet.DataPacket, error) {
		return nil, nil
	})
	return b.server.Serve()
}

func (b *broker) Shutdown() error {
	return b.server.Shutdown(context.Background())
}

func (b *broker) SendData(ctx context.Context, channel string, payload []byte) error {
	if users, ok := b.channels.Get(channel); ok {
		users.Range(func(uid string) bool {
			b.taskQueue.Push(func() {
				if conn, err := b.server.GetConn(uid); err != nil {
					conn.SendData(payload)
				}
			})
			return true
		})
	}

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

	return nil
}

func (b *broker) Leave(channels []string, uid string) error {
	for _, channel := range channels {
		if users, ok := b.channels.Get(channel); ok {
			users.Del(uid)
		}
	}

	return nil
}
