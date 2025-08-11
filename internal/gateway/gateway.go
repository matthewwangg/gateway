package gateway

import (
	"context"
	"net/http"
)

type Gateway struct {
	Server *http.Server
}

func NewGateway() *Gateway {
	mux := http.NewServeMux()
	SetupRoutes(mux)
	return &Gateway{
		Server: &http.Server{
			Addr:    ":8080",
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
