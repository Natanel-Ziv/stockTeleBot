package infrastructure

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var AppContext context.Context

func init() {
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()

		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

		<-signalCh
	}()

	AppContext = ctxWithCancel
}