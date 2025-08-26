package gateway

import (
	"encoding/json"
	"net/http"
	"os"

	client "github.com/matthewwangg/gateway/internal/client"
	middleware "github.com/matthewwangg/gateway/internal/middleware"
	models "github.com/matthewwangg/gateway/internal/models"
)

func (g *Gateway) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (g *Gateway) Login(w http.ResponseWriter, r *http.Request) {
	type LoginRequestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var body LoginRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Username != os.Getenv("GATEWAY_USER") || body.Password != os.Getenv("GATEWAY_PASSWORD") {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	type LoginResponseBody struct {
		Token string `json:"token"`
	}

	token, err := middleware.GenerateJWTToken(body.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := LoginResponseBody{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (g *Gateway) Reload(w http.ResponseWriter, r *http.Request) {
	g.Registry.Reload()
}

func (g *Gateway) Services(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(g.Registry.Services); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Type != models.REST {
		http.Error(w, "api type not yet supported", http.StatusBadRequest)
		return
	}

	service, ok := g.Registry.Services[body.Service]
	if !ok {
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}

	c := client.NewRESTClient(service)
	if c == nil {
		http.Error(w, "service not healthy", http.StatusNotFound)
		return
	}

	result := map[string]interface{}{}
	if err := c.Call(body.Endpoint, body.Params, &result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
