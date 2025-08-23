package gateway

import (
	"net/http"
)

func (g *Gateway) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (g *Gateway) Reload(w http.ResponseWriter, r *http.Request) {
	g.Registry.Reload()
}
