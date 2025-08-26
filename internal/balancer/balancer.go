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
	Address      string
	ResponseTime time.Duration
}

type ServiceUsage struct {
	AddressUsages []*AddressUsage
	Count         int
}

type LoadBalancer struct {
	Mode          LoadBalancerMode
	ServiceUsages map[string]*ServiceUsage
}

func NewLoadBalancer(mode LoadBalancerMode, services map[string]*models.ServiceDefinition) *LoadBalancer {
	serviceUsages := make(map[string]*ServiceUsage)
	for serviceName, serviceDefinition := range services {
		addressUsages := make([]*AddressUsage, 0)
		for _, address := range serviceDefinition.Addresses {
			addressUsage := &AddressUsage{
				Address:      address,
				ResponseTime: 0 * time.Second,
			}
			addressUsages = append(addressUsages, addressUsage)
		}

		serviceUsages[serviceName] = &ServiceUsage{
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
		address = usage.AddressUsages[(usage.Count)%len(usage.AddressUsages)].Address
		usage.Count += 1
	case LeastResponseTime:
		lowestResponseTime := 24 * time.Hour
		for _, addressUsage := range usage.AddressUsages {
			if addressUsage.ResponseTime < lowestResponseTime {
				lowestResponseTime = addressUsage.ResponseTime
				address = addressUsage.Address
			}
		}
	}

	return address, nil
}
