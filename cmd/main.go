package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	gateway "github.com/matthewwangg/gateway/internal/gateway"
)

func main() {
	g := gateway.NewGateway()

	go func() {
		log.Println("gateway listening on port 8080")
		if err := g.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start gateway: %s", err)
		}
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("gateway shutting down")
	if err := g.Stop(ctx); err != nil {
		log.Fatalf("failed to stop gateway: %s", err)
	}
}
