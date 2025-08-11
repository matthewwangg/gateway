package gateway

import (
	"context"
	"net/http"
)

type Gateway struct {
	Server *http.Server
}

func NewGateway() *Gateway {
	return &Gateway{
		Server: &http.Server{
			Addr:    ":8080",
			Handler: http.NewServeMux(),
		},
	}
}

func (g *Gateway) Start() error {
	return g.Server.ListenAndServe()
}

func (g *Gateway) Stop(ctx context.Context) error {
	return g.Server.Shutdown(ctx)
}
