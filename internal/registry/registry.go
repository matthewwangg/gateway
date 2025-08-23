package registry

import (
	models "github.com/matthewwangg/gateway/internal/models"
	"github.com/matthewwangg/gateway/internal/parser"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ServiceRegistry struct {
	Services  map[string]*models.ServiceDefinition
	Directory string
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
		log.Fatalf("error reading directory: %v", err)
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

	s.Services = services
}

func (s *ServiceRegistry) Register(service *models.ServiceDefinition) {
	s.Services[service.Name] = service
}

func (s *ServiceRegistry) Get(serviceName string) *models.ServiceDefinition {
	if service, ok := s.Services[serviceName]; ok {
		return service
	}
	return nil
}
