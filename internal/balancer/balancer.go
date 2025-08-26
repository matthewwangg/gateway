package balancer

import (
	"errors"
	"time"

	models "github.com/matthewwangg/gateway/internal/models"
)

type LoadBalancerMode string

const (
	RoundRobin        LoadBalancerMode = "round_robin"
	LeastResponseTime LoadBalancerMode = "least_response_time"
)

type AddressUsage struct {
	ResponseTime time.Duration
}

type ServiceUsage struct {
	Addresses     []string
	AddressUsages map[string]*AddressUsage
	Count         int
}

type LoadBalancer struct {
	Mode          LoadBalancerMode
	ServiceUsages map[string]*ServiceUsage
}

func NewLoadBalancer(mode LoadBalancerMode, services map[string]*models.ServiceDefinition) *LoadBalancer {
	serviceUsages := make(map[string]*ServiceUsage)
	for serviceName, serviceDefinition := range services {
		addressUsages := make(map[string]*AddressUsage)
		for _, address := range serviceDefinition.Addresses {
			addressUsages[address] = &AddressUsage{
				ResponseTime: 0 * time.Second,
			}
		}

		serviceUsages[serviceName] = &ServiceUsage{
			Addresses:     serviceDefinition.Addresses,
			AddressUsages: addressUsages,
			Count:         0,
		}
	}

	return &LoadBalancer{
		Mode:          mode,
		ServiceUsages: serviceUsages,
	}
}

func (l *LoadBalancer) Select(serviceName string) (string, error) {
	usage, ok := l.ServiceUsages[serviceName]
	if !ok || len(usage.AddressUsages) < 1 {
		return "", errors.New("service not valid")
	}

	address := ""

	switch l.Mode {
	case RoundRobin:
		address = usage.Addresses[(usage.Count)%len(usage.AddressUsages)]
		usage.Count += 1
	case LeastResponseTime:
		lowestResponseTime := 24 * time.Hour
		for addressKey, addressUsage := range usage.AddressUsages {
			if addressUsage.ResponseTime < lowestResponseTime {
				lowestResponseTime = addressUsage.ResponseTime
				address = addressKey
			}
		}
	}

	return address, nil
}
