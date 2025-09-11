package registry

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	logger "github.com/matthewwangg/gateway/internal/logger"
	models "github.com/matthewwangg/gateway/internal/models"
	parser "github.com/matthewwangg/gateway/internal/parser"
)

type ServiceRegistry struct {
	Services  map[string]*models.ServiceDefinition
	Directory string
	mu        sync.RWMutex
}

func NewServiceRegistry(directory string) *ServiceRegistry {
	serviceRegistry := &ServiceRegistry{
		Services:  map[string]*models.ServiceDefinition{},
		Directory: directory,
	}

	serviceRegistry.Reload()

	return serviceRegistry
}

func (s *ServiceRegistry) Reload() {
	files, err := os.ReadDir(s.Directory)
	if err != nil {
		logger.Log.Error("error reading directory: " + string(err.Error()))
		return
	}

	services := map[string]*models.ServiceDefinition{}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".svc") {
			continue
		}

		path := filepath.Join(s.Directory, file.Name())
		fileParser := parser.NewParser(path)

		serviceDefinition := fileParser.Parse()
		services[serviceDefinition.Name] = serviceDefinition
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.Services = services
}

func (s *ServiceRegistry) Register(service *models.ServiceDefinition) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Services[service.Name] = service
}

func (s *ServiceRegistry) Get(serviceName string) *models.ServiceDefinition {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if service, ok := s.Services[serviceName]; ok {
		return service
	}
	return nil
}
