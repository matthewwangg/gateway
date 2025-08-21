package models

type APIType string

const (
	gRPC APIType = "grpc"
	REST APIType = "rest"
)

type ServiceDefinition struct {
	Name      string
	Replicas  int
	Addresses []string
	APIType   APIType
	Endpoints []string
}
