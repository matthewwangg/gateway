package balancer

import (
	"errors"
	
	models "github.com/matthewwangg/gateway/internal/models"
)

type LoadBalancerMode string

const (
	RoundRobin LoadBalancerMode = "round_robin"
)

type ServiceUsage struct {
	Addresses []string
	Count     int
}

type LoadBalancer struct {
	Mode          LoadBalancerMode
	ServiceUsages map[string]*ServiceUsage
}

func NewLoadBalancer(mode LoadBalancerMode, services map[string]*models.ServiceDefinition) *LoadBalancer {
	serviceUsages := make(map[string]*ServiceUsage)

	for serviceName, serviceDefinition := range services {
		serviceUsages[serviceName] = &ServiceUsage{
			Addresses: serviceDefinition.Addresses,
			Count:     0,
		}
	}

	return &LoadBalancer{
		Mode:          mode,
		ServiceUsages: serviceUsages,
	}
}

func (l *LoadBalancer) Select(serviceName string) (string, error) {
	usage, ok := l.ServiceUsages[serviceName]
	if !ok || len(usage.Addresses) < 1 {
		return "", errors.New("service not valid")
	}

	address := ""

	switch l.Mode {
	case RoundRobin:
		address = usage.Addresses[(usage.Count)%len(usage.Addresses)]
		usage.Count += 1
	}

	return address, nil
}
