package main

import (
	"context"

	"sutext.github.io/entry/logger"
	"sutext.github.io/entry/packet"
	"sutext.github.io/entry/server"
)

func main() {
	config := server.NewConfig()
	config.LoggerLevel = logger.LevelDebug
	ctx := context.Background()
	s := server.New(config)

	s.HandleLogin(func(p *packet.Identity) error {
		return nil
	})
	s.HandleData(func(p *packet.DataPacket) (*packet.DataPacket, error) {
		return nil, nil
	})
	s.Listen(ctx)
}
