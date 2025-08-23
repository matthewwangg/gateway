package client

import (
	"github.com/matthewwangg/gateway/internal/models"
	"net/http"
)

type RESTClient struct {
	Address   string
	Endpoints []RESTEndpoint
}

type RESTEndpoint struct {
	Path   string
	Method RESTMethod
}

type RESTMethod string

const (
	GET    RESTMethod = "GET"
	POST   RESTMethod = "POST"
	PUT    RESTMethod = "PUT"
	PATCH  RESTMethod = "PATCH"
	DELETE RESTMethod = "DELETE"
)

func NewRESTClient(address string, serviceDefinition *models.ServiceDefinition) *RESTClient {
	client := &RESTClient{
		Address:   address,
		Endpoints: []RESTEndpoint{},
	}

	if !client.HealthCheck() {
		return nil
	}

	return client
}

func (c *RESTClient) HealthCheck() bool {
	response, err := http.Get(c.Address + "/healthz")
	if err != nil {
		return false
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			return
		}
	}()

	return response.StatusCode == 200
}
