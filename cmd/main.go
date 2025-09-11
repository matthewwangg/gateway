package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"

	gateway "github.com/matthewwangg/gateway/internal/gateway"
	logger "github.com/matthewwangg/gateway/internal/logger"
)

func main() {
	logger.Init(logger.Local, "gateway", "127.0.0.1")

	g := gateway.NewGateway()

	go func() {
		logger.Log.Info("gateway listening on port 8080")
		if err := g.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Error("failed to start gateway: " + string(err.Error()))
			os.Exit(1)
		}
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Log.Info("gateway shutting down")
	if err := g.Stop(ctx); err != nil {
		logger.Log.Error("failed to stop gateway: " + string(err.Error()))
		os.Exit(1)
	}
}
