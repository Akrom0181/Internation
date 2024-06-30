package client

import (
	"fmt"
	"user_service/config"
	"user_service/genproto/schedule_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceManagerI interface {
	ScheduleService() schedule_service.ScheduleServiceClient
}

type grpcClients struct {
	scheduleService schedule_service.ScheduleServiceClient
}

func NewGrpcClients(cfg config.Config) (ServiceManagerI, error) {

	connScheduleService, err := grpc.Dial(
		fmt.Sprintf("%s%s", cfg.ScheduleServicePort, cfg.ScheduleServiceHost),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(52428800), grpc.MaxCallSendMsgSize(52428800)),
	)
	if err != nil {
		return nil, err
	}

	return &grpcClients{
		scheduleService: schedule_service.NewScheduleServiceClient(connScheduleService),
	}, nil
}

func (g *grpcClients) ScheduleService() schedule_service.ScheduleServiceClient {
	return g.scheduleService
}
