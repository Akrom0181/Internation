package client

import "schedule_service/config"

type ServiceManagerI interface{}

type grpcClients struct{}

func NewGrpcClients(cfg config.Config) (ServiceManagerI, error) {

	return &grpcClients{}, nil
}
