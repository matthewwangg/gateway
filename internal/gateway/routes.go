package gateway

import (
	"net/http"
)

func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", HealthCheck)
}
