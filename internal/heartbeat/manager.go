package heartbeat

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/matthewwangg/gateway/internal/logger"
	"github.com/matthewwangg/gateway/internal/models"
)

type Manager struct {
	Services map[string]*models.ServiceDefinition
	Interval time.Duration
	Timeout  time.Duration
	Context  context.Context
	Cancel   context.CancelFunc
}

func NewManager(services map[string]*models.ServiceDefinition, interval time.Duration, timeout time.Duration) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		Services: services,
		Interval: interval,
		Timeout:  timeout,
		Context:  ctx,
		Cancel:   cancel,
	}
}

func (m *Manager) Start() {
	ticker := time.NewTicker(m.Interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				m.checkServices()
			case <-m.Context.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (m *Manager) Stop() {
	m.Cancel()
}

func (m *Manager) checkServices() {
	var wg sync.WaitGroup
	for _, service := range m.Services {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.checkService(service)
		}()
	}
}

func (m *Manager) checkService(service *models.ServiceDefinition) {
	healthEndpoint := service.HealthEndpoint
	for _, address := range service.Addresses {
		healthy := m.checkAddress(address, healthEndpoint)
		if !healthy {
			logger.Log.Warn("unhealthy address found for " + service.Name)
		}
	}
}

func (m *Manager) checkAddress(address string, healthEndpoint string) bool {
	url := fmt.Sprintf("http://%s%s", address, healthEndpoint)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return false
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			logger.Log.Warn(err.Error())
		}
	}()

	return response.StatusCode == http.StatusOK
}
