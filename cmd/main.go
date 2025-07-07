package main

import (
	"context"

	"sutext.github.io/entry/packet"
	"sutext.github.io/entry/server"
)

func main() {
	ctx := context.Background()
	s := server.New()

	s.HandleAuth(func(p *packet.Identity) error {
		return nil
	})
	s.HandleData(func(p *packet.DataPacket) (*packet.DataPacket, error) {
		return nil, nil
	})
	s.Listen(ctx, "8080")
}
