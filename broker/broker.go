package broker

import (
	"context"

	"github.com/redis/go-redis/v9"
	"sutext.github.io/entry/nio"
)

type Broker interface {
	Publish(ctx context.Context, topic string, payload []byte) error
	Subscribe(uid int64, topics []string) error
	Unsubscribe(uid int64, topics []string) error
}
type broker struct {
	config *Config
	server nio.Server
	rdb    *redis.Client
}

func New() Broker {
	conf := DefaultConfig()
	return &broker{
		config: conf,
	}
}

func (b *broker) Run() error {
	b.server = nio.New(b.config.Port)
	b.rdb = redis.NewClient(&redis.Options{
		Addr:     b.config.Redis.Addr,
		Password: b.config.Redis.Password,
		DB:       b.config.Redis.DB,
	})
	return b.server.Start()
}
func (b *broker) Shutdown() error {
	return b.server.Shutdown(nil)
}

func (b *broker) Publish(ctx context.Context, topic string, payload []byte) error {

	// cmd := b.rdb.SScan(ctx,b.topicKey(topic),0,"",0).Err()
	return nil
}
func (b *broker) Subscribe(uid int64, topics []string) error {
	ctx := context.Background()
	for _, topic := range topics {
		b.rdb.SAdd(ctx, b.topicKey(topic), uid)
	}
	return nil
}
func (b *broker) Unsubscribe(uid int64, topics []string) error {
	ctx := context.Background()
	for _, topic := range topics {
		b.rdb.SRem(ctx, b.topicKey(topic), uid)
	}
	return nil
}
