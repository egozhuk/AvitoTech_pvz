package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"pvs/internal/app"
	"pvs/internal/config"
	"syscall"
)

func main() {
	_ = config.Load()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to start app: %v", err)
	}

	go func() {
		if err := a.Run(); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	if err := a.Stop(ctx); err != nil {
		log.Fatalf("failed to shutdown: %v", err)
	}
}
