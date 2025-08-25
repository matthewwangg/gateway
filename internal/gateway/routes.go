package gateway

import (
	"net/http"
)

func (g *Gateway) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", g.HealthCheck)
	mux.HandleFunc("GET /reload", g.Reload)
	mux.HandleFunc("GET /services", g.Services)
	mux.HandleFunc("POST /call", g.Call)
}
