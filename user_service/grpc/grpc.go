package grpc

import (
	"user_service/config"
	"user_service/genproto/user_service"
	"user_service/grpc/client"
	"user_service/grpc/service"
	"user_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvc client.ServiceManagerI) (grpcServer *grpc.Server) {

	grpcServer = grpc.NewServer()

	user_service.RegisterAdministrationServiceServer(grpcServer, service.NewAdministrationService(cfg, log, strg, srvc))
	user_service.RegisterBranchServiceServer(grpcServer, service.NewBranchService(cfg, log, strg, srvc))
	user_service.RegisterManagerServiceServer(grpcServer, service.NewManagerService(cfg, log, strg, srvc))
	user_service.RegisterStudentServiceServer(grpcServer, service.NewStudentService(cfg, log, strg, srvc))
	user_service.RegisterSupportTeacherServiceServer(grpcServer, service.NewSupportTeacherService(cfg, log, strg, srvc))
	user_service.RegisterTeacherServiceServer(grpcServer, service.NewTeacherService(cfg, log, strg, srvc))
	user_service.RegisterLoginServiceServer(grpcServer, service.NewLoginService(cfg, log, strg, srvc))
	reflection.Register(grpcServer)
	return
}
