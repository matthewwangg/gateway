package gateway

import (
	"context"
	"net/http"
	"os"
	"time"

	balancer "github.com/matthewwangg/gateway/internal/balancer"
	heartbeat "github.com/matthewwangg/gateway/internal/heartbeat"
	metrics "github.com/matthewwangg/gateway/internal/metrics"
	middleware "github.com/matthewwangg/gateway/internal/middleware"
	registry "github.com/matthewwangg/gateway/internal/registry"
)

type Gateway struct {
	Server           *http.Server
	Registry         *registry.ServiceRegistry
	LoadBalancer     *balancer.LoadBalancer
	HeartbeatManager *heartbeat.Manager
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
	heartbeatManager := heartbeat.NewManager(serviceRegistry.Services, time.Minute, 10*time.Second)

	g := &Gateway{
		Server:           server,
		Registry:         serviceRegistry,
		LoadBalancer:     loadBalancer,
		HeartbeatManager: heartbeatManager,
	}

	mux := http.NewServeMux()
	g.SetupRoutes(mux)

	handler := middleware.Use(mux,
		middleware.Logger,
		middleware.RateLimiter,
	)
	g.Server.Handler = handler

	metrics.Init()

	return g
}

func (g *Gateway) Start() error {
	g.HeartbeatManager.Start()
	return g.Server.ListenAndServe()
}

func (g *Gateway) Stop(ctx context.Context) error {
	g.HeartbeatManager.Stop()
	return g.Server.Shutdown(ctx)
}
