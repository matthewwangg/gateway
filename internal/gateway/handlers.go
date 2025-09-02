package gateway

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	balancer "github.com/matthewwangg/gateway/internal/balancer"
	client "github.com/matthewwangg/gateway/internal/client"
	metrics "github.com/matthewwangg/gateway/internal/metrics"
	middleware "github.com/matthewwangg/gateway/internal/middleware"
	models "github.com/matthewwangg/gateway/internal/models"
)

func (g *Gateway) HealthCheck(w http.ResponseWriter, r *http.Request) {
	metrics.Tracker.RecordRequest("/healthz", 200)
	w.WriteHeader(http.StatusOK)
}

func (g *Gateway) Metrics(w http.ResponseWriter, r *http.Request) {
	metrics.Tracker.RecordRequest("/metrics", 200)
	promhttp.Handler().ServeHTTP(w, r)
}

func (g *Gateway) Login(w http.ResponseWriter, r *http.Request) {
	type LoginRequestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var body LoginRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		metrics.Tracker.RecordRequest("/login", 400)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Username != os.Getenv("GATEWAY_USER") || body.Password != os.Getenv("GATEWAY_PASSWORD") {
		metrics.Tracker.RecordRequest("/login", 401)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	type LoginResponseBody struct {
		Token string `json:"token"`
	}

	token, err := middleware.GenerateJWTToken(body.Username)
	if err != nil {
		metrics.Tracker.RecordRequest("/login", 500)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := LoginResponseBody{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		metrics.Tracker.RecordRequest("/login", 500)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	metrics.Tracker.RecordRequest("/login", 200)
}

func (g *Gateway) Reload(w http.ResponseWriter, r *http.Request) {
	metrics.Tracker.RecordRequest("/reload", 200)
	g.Registry.Reload()
	g.LoadBalancer = balancer.NewLoadBalancer(g.LoadBalancer.Mode, g.Registry.Services)
}

func (g *Gateway) Services(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(g.Registry.Services); err != nil {
		metrics.Tracker.RecordRequest("/services", 500)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	metrics.Tracker.RecordRequest("/services", 200)
}

func (g *Gateway) Call(w http.ResponseWriter, r *http.Request) {
	type CallRequestBody struct {
		Type     models.APIType         `json:"type"`
		Service  string                 `json:"service"`
		Endpoint string                 `json:"endpoint"`
		Params   map[string]interface{} `json:"params"`
	}

	var body CallRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		metrics.Tracker.RecordRequest("/call", 400)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Type != models.REST {
		metrics.Tracker.RecordRequest("/call", 400)
		http.Error(w, "api type not yet supported", http.StatusBadRequest)
		return
	}

	service, ok := g.Registry.Services[body.Service]
	if !ok {
		metrics.Tracker.RecordRequest("/call", 404)
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}

	c := client.NewRESTClient(service, g.LoadBalancer)
	if c == nil {
		metrics.Tracker.RecordRequest("/call", 404)
		http.Error(w, "service not healthy", http.StatusNotFound)
		return
	}

	result := map[string]interface{}{}
	start := time.Now()
	if err := c.Call(body.Endpoint, body.Params, &result); err != nil {
		metrics.Tracker.RecordRequest("/call", 500)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	duration := time.Since(start)
	g.LoadBalancer.ServiceUsages[body.Service].AddressUsages[c.Address].ResponseTime = duration

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		metrics.Tracker.RecordRequest("/call", 500)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	metrics.Tracker.RecordServiceCall(body.Service)
}
