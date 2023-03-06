//go:build linux || darwin
// +build linux darwin

package proc

import (
	"context"
	"github.com/mangohow/httputil/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var once sync.Once

func SetupSignalHandler(log logger.Logger) context.Context {
	var c context.Context

	once.Do(func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2)

		ctx, cancel := context.WithCancel(context.Background())
		c = ctx

		go func() {
			setLogger(log)
			var stopper Stopper
			exit := false
			for {
				sig := <-ch

				switch sig {
				case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
					log.Infof("receive signal %s", sig.String())
					if exit {
						os.Exit(0)
					}
					cancel()
					exit = true
					return
				case syscall.SIGUSR1:
					if stopper == nil {
						stopper = StartProfile()
					} else {
						stopper.Stop()
						stopper = nil
					}
				case syscall.SIGUSR2:
					dumpGoroutines()
				}

			}
		}()
	})

	return c
}
