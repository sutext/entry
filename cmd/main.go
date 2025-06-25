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

	s.HandleLogin(func(ctx context.Context, conn *server.Conn, p packet.ConnectPacket) error {
		return nil
	})
	s.HandleData(func(ctx context.Context, conn *server.Conn, p packet.DataPacket) error {
		return nil
	})
	s.Run(ctx)
}
