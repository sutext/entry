package entry

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sutext.github.io/entry/server"
)

type Entry struct {
	logger *slog.Logger
}

func (e *Entry) ListenAndServe(ctx context.Context) error {
	e.logger.InfoContext(ctx, "entry server start")
	ctx, cancel := context.WithCancelCause(ctx)
	server := server.New()
	done := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel(fmt.Errorf("entry signal received"))
	}()
	go func() {
		defer close(done)
		<-ctx.Done()
		server.Shutdown(ctx)
		// s.tcpListener.shutdown(ctx)
		// s.peerListener.shutdown(ctx)
	}()

	// go s.tcpListener.listen(ctx)
	// go s.peerListener.listen(ctx)
	<-ctx.Done()
	e.logger.InfoContext(ctx, "entry server stoped")
	timeout := time.NewTimer(time.Second * 15)
	defer timeout.Stop()
	select {
	case <-timeout.C:
		e.logger.WarnContext(ctx, "entry server graceful shutdown timeout")
	case <-done:
		e.logger.DebugContext(ctx, "entry server graceful shutdown")
	}
	return context.Cause(ctx)
}
