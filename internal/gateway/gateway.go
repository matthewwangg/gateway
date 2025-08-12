package gateway

import (
	"context"
	"net/http"
	"os"
)

type Gateway struct {
	Server *http.Server
}

func NewGateway() *Gateway {
	addr := os.Getenv("GATEWAY_ADDRESS")
	if addr == "" {
		addr = ":8080"
	}

	mux := http.NewServeMux()
	SetupRoutes(mux)

	return &Gateway{
		Server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func (g *Gateway) Start() error {
	return g.Server.ListenAndServe()
}

func (g *Gateway) Stop(ctx context.Context) error {
	return g.Server.Shutdown(ctx)
}
