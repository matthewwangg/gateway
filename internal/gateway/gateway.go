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

	mux := http.NewServeMux()
	SetupRoutes(mux)

	return &Gateway{
		Server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		Registry: registry.NewServiceRegistry(os.Getenv("REGISTRY_DIRECTORY")),
	}
}

func (g *Gateway) Start() error {
	return g.Server.ListenAndServe()
}

func (g *Gateway) Stop(ctx context.Context) error {
	return g.Server.Shutdown(ctx)
}
