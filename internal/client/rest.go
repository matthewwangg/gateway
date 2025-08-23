package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	models "github.com/matthewwangg/gateway/internal/models"
)

type RESTClient struct {
	Address   string
	Endpoints map[string]RESTEndpoint
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
		endpoints := make(map[string]RESTEndpoint)
		for _, endpoint := range serviceDefinition.Endpoints {
			parts := strings.Split(endpoint, " ")

			if len(parts) != 2 {
				endpoints["GET "+endpoint] = RESTEndpoint{
					Path:   endpoint,
					Method: GET,
				}
				endpoints["POST "+endpoint] = RESTEndpoint{
					Path:   endpoint,
					Method: POST,
				}
				endpoints["PUT "+endpoint] = RESTEndpoint{
					Path:   endpoint,
					Method: PUT,
				}
				endpoints["PATCH "+endpoint] = RESTEndpoint{
					Path:   endpoint,
					Method: PATCH,
				}
				endpoints["DELETE "+endpoint] = RESTEndpoint{
					Path:   endpoint,
					Method: DELETE,
				}
			} else {
				method := RESTMethod(parts[0])
				path := parts[1]
				endpoints[endpoint] = RESTEndpoint{
					Path:   path,
					Method: method,
				}
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

func (c *RESTClient) Call(endpointName string, params map[string]interface{}, result interface{}) error {
	address := c.Address

	if !strings.HasPrefix(address, "http://") {
		address = "http://" + address
	}

	endpoint, ok := c.Endpoints[endpointName]
	if !ok {
		return fmt.Errorf("endpoint %s not found", endpointName)
	}

	url := address + endpoint.Path

	var request *http.Request
	var err error

	if endpoint.Method == GET {
		request, err = http.NewRequest(string(endpoint.Method), url, nil)
		if err != nil {
			return err
		}
	} else {
		body, err := json.Marshal(params)
		if err != nil {
			return err
		}

		request, err = http.NewRequest(string(endpoint.Method), url, bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		if err = response.Body.Close(); err != nil {
			return
		}
	}()

	if result == nil {
		return nil
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &result); err != nil {
		return err
	}

	return nil
}
