package gateway

import (
	"context"
	"net/http"
	"os"

	registry "github.com/matthewwangg/gateway/internal/registry"
)

type Gateway struct {
	Server   *http.Server
	Registry *registry.ServiceRegistry
}

func NewGateway() *Gateway {
	addr := os.Getenv("GATEWAY_ADDRESS")
	if addr == "" {
		addr = ":8080"
	}

	g := &Gateway{
		Server: &http.Server{
			Addr:    addr,
			Handler: nil,
		},
		Registry: registry.NewServiceRegistry(os.Getenv("REGISTRY_DIRECTORY")),
	}

	mux := http.NewServeMux()
	g.SetupRoutes(mux)
	g.Server.Handler = mux

	return g
}

func (g *Gateway) Start() error {
	return g.Server.ListenAndServe()
}

func (g *Gateway) Stop(ctx context.Context) error {
	return g.Server.Shutdown(ctx)
}
