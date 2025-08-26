package gateway

import (
	"context"
	"net/http"
	"os"

	balancer "github.com/matthewwangg/gateway/internal/balancer"
	registry "github.com/matthewwangg/gateway/internal/registry"
)

type Gateway struct {
	Server       *http.Server
	Registry     *registry.ServiceRegistry
	LoadBalancer *balancer.LoadBalancer
}

func NewGateway() *Gateway {
	addr := os.Getenv("GATEWAY_ADDRESS")
	if addr == "" {
		addr = ":8080"
	}

	server := &http.Server{
		Addr:    addr,
		Handler: nil,
	}

	serviceRegistry := registry.NewServiceRegistry(os.Getenv("REGISTRY_DIRECTORY"))
	loadBalancer := balancer.NewLoadBalancer(balancer.RoundRobin, serviceRegistry.Services)

	g := &Gateway{
		Server:       server,
		Registry:     serviceRegistry,
		LoadBalancer: loadBalancer,
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
