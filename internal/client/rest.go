package client

import (
	"net/http"
	"strings"

	models "github.com/matthewwangg/gateway/internal/models"
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

func NewRESTClient(serviceDefinition *models.ServiceDefinition) *RESTClient {
	for _, address := range serviceDefinition.Addresses {
		endpoints := make([]RESTEndpoint, 0)
		for _, endpoint := range serviceDefinition.Endpoints {
			parts := strings.Split(endpoint, " ")

			if len(parts) != 2 {
				endpoints = append(endpoints, RESTEndpoint{
					Path:   endpoint,
					Method: GET,
				})
			} else {
				method := RESTMethod(parts[0])
				path := parts[1]
				endpoints = append(endpoints, RESTEndpoint{
					Path:   path,
					Method: method,
				})
			}
		}

		client := &RESTClient{
			Address:   address,
			Endpoints: endpoints,
		}

		if client.HealthCheck() {
			return client
		}
	}

	return nil
}

func (c *RESTClient) HealthCheck() bool {
	response, err := http.Get("http://" + c.Address + "/healthz")
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
