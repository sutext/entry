package entry

import (
	"context"

	"sutext.github.io/entry/server"
)

func Start(config *server.Config) {
	ctx := context.Background()
	server := server.New(config)
	server.Run(ctx)
}
