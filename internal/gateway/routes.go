package gateway

import (
	"net/http"

	middleware "github.com/matthewwangg/gateway/internal/middleware"
)

func (g *Gateway) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", g.HealthCheck)
	mux.HandleFunc("POST /login", g.Login)

	mux.Handle("GET /reload", g.Secure(g.Reload))
	mux.Handle("GET /services", g.Secure(g.Services))
	mux.Handle("POST /call", g.Secure(g.Call))
}

func (g *Gateway) Secure(handler http.HandlerFunc) http.Handler {
	return middleware.Use(handler, middleware.CheckJWTToken)
}
