package grpc

import (
	"schedule_service/config"
	"schedule_service/genproto/schedule_service"
	"schedule_service/grpc/client"
	"schedule_service/grpc/service"
	"schedule_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvc client.ServiceManagerI) (grpcServer *grpc.Server) {

	grpcServer = grpc.NewServer()

	schedule_service.RegisterEventStudentServiceServer(grpcServer, service.NewEventStudentService(cfg, log, strg, srvc))
	schedule_service.RegisterEventServiceServer(grpcServer, service.NewEventService(cfg, log, strg, srvc))
	schedule_service.RegisterGroupServiceServer(grpcServer, service.NewGroupService(cfg, log, strg, srvc))
	schedule_service.RegisterJournalServiceServer(grpcServer, service.NewJournalService(cfg, log, strg, srvc))
	schedule_service.RegisterScheduleServiceServer(grpcServer, service.NewScheduleService(cfg, log, strg, srvc))
	schedule_service.RegisterStudentPaymentServiceServer(grpcServer, service.NewStudentPaymentService(cfg, log, strg, srvc))
	schedule_service.RegisterStudentTaskServiceServer(grpcServer, service.NewStudentTaskService(cfg, log, strg, srvc))
	schedule_service.RegisterTaskServiceServer(grpcServer, service.NewTaskService(cfg, log, strg, srvc))
	reflection.Register(grpcServer)
	return
}
