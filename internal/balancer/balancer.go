package balancer

type LoadBalancerMode string

const (
	LoadBalancerModeRoundRobin LoadBalancerMode = "round_robin"
)

type LoadBalancer struct {
	Mode LoadBalancerMode
}

func NewLoadBalancer(mode LoadBalancerMode) *LoadBalancer {
	return &LoadBalancer{mode}
}
