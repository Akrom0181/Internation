package service

import (
	"context"
	"user_service/config"
	"user_service/genproto/user_service"
	"user_service/grpc/client"

	"user_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
)

type ManagerService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
}

func NewManagerService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *ManagerService {
	return &ManagerService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (s *ManagerService) Create(ctx context.Context, req *user_service.CreateManager) (*user_service.Manager, error) {
	s.log.Info("---CreateManager--->>>", logger.Any("req", req))

	resp, err := s.strg.Manager().Create(ctx, req)
	if err != nil {
		s.log.Error("---CreateManager--->>>", logger.Error(err))
		return &user_service.Manager{}, err
	}

	return resp, nil
}

func (s *ManagerService) GetByID(ctx context.Context, req *user_service.ManagerPrimaryKey) (*user_service.Manager, error) {
	s.log.Info("---GetSingleManager--->>>", logger.Any("req", req))

	resp, err := s.strg.Manager().GetByID(ctx, req)
	if err != nil {
		s.log.Error("---GetSingleManager--->>>", logger.Error(err))
		return &user_service.Manager{}, err
	}

	return resp, nil
}

func (s *ManagerService) GetList(ctx context.Context, req *user_service.GetListManagerRequest) (*user_service.GetListManagerResponse, error) {
	s.log.Info("---GetAllManagers--->>>", logger.Any("req", req))

	resp, err := s.strg.Manager().GetList(ctx, req)
	if err != nil {
		s.log.Error("---GetAllManagers--->>>", logger.Error(err))
		return &user_service.GetListManagerResponse{}, err
	}

	return resp, nil
}

func (s *ManagerService) Update(ctx context.Context, req *user_service.UpdateManager) (*user_service.Manager, error) {
	s.log.Info("---UpdateManager--->>>", logger.Any("req", req))

	resp, err := s.strg.Manager().Update(ctx, req)
	if err != nil {
		s.log.Error("---UpdateManager--->>>", logger.Error(err))
		return &user_service.Manager{}, err
	}

	return resp, nil
}

func (s *ManagerService) Delete(ctx context.Context, req *user_service.ManagerPrimaryKey) (*user_service.EmptyManager, error) {
	s.log.Info("---DeleteManager--->>>", logger.Any("req", req))

	err := s.strg.Manager().Delete(ctx, req)
	if err != nil {
		s.log.Error("---DeleteManager--->>>", logger.Error(err))
		return &user_service.EmptyManager{}, err
	}

	return &user_service.EmptyManager{}, nil
}
