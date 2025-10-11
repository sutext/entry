package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sutext.github.io/entry/logger"
	"sutext.github.io/entry/nio"
	"sutext.github.io/entry/packet"
	"sutext.github.io/entry/server"
)

func main() {
	runNIO()
	// runServer()
}
func runServer() {
	ctx := context.Background()
	logger := logger.New(logger.LevelDebug, logger.FormatJSON)
	s := server.New()
	s.HandleAuth(func(p *packet.Identity) error {
		return nil
	})
	s.HandleData(func(p *packet.DataPacket) (*packet.DataPacket, error) {
		return nil, nil
	})
	s.SetLogger(logger)
	s.SetKeepAlive(time.Second*60, time.Second*5)
	logger.InfoContext(ctx, "entry server start")
	ctx, cancel := context.WithCancelCause(ctx)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		logger.InfoContext(ctx, "entry signal received")
		cancel(fmt.Errorf("entry signal received"))
	}()
	go s.Listen(ctx, "8080")
	<-ctx.Done()
	logger.InfoContext(ctx, "entry server start graceful shutdown")
	done := make(chan struct{})
	go func() {
		s.Shutdown(ctx)
		close(done)
	}()
	timeout := time.NewTimer(time.Second * 15)
	defer timeout.Stop()
	select {
	case <-timeout.C:
		logger.WarnContext(ctx, "entry server graceful shutdown timeout")
	case <-done:
		logger.DebugContext(ctx, "entry server graceful shutdown")
	}
}
func runNIO() {
	s := nio.New("8080")
	s.HandleAuth(func(p *packet.Identity) error {
		return nil
	})
	s.HandleData(func(uid string, p *packet.DataPacket) (*packet.DataPacket, error) {
		return nil, nil
	})
	s.Start()
}
