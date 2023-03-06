//go:build windows
// +build windows

package proc

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var onlyOneSignalShot = make(chan struct{})

func SetupSignalHandler() context.Context {
	close(onlyOneSignalShot)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		cancel()
		<-c
		os.Exit(0)
	}()

	return ctx
}
