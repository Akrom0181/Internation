package service

import (
	"context"
	"user_service/config"
	"user_service/genproto/user_service"
	"user_service/grpc/client"
	"user_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
)

type UserService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
}

func NewAdministrationService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *UserService {
	return &UserService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (f *UserService) Create(ctx context.Context, req *user_service.CreateAdministration) (*user_service.Administration, error) {

	f.log.Info("---CreateAdministration--->>>", logger.Any("req", req))

	resp, err := f.strg.Administration().Create(ctx, req)
	if err != nil {
		f.log.Error("---CreateAdministration--->>>", logger.Error(err))
		return &user_service.Administration{}, err
	}

	return resp, nil
}

func (f *UserService) GetByID(ctx context.Context, req *user_service.AdministrationPrimaryKey) (*user_service.Administration, error) {
	f.log.Info("---GetSingleAdministration--->>>", logger.Any("req", req))

	resp, err := f.strg.Administration().GetByID(ctx, req)
	if err != nil {
		f.log.Error("---GetSingleAdministration--->>>", logger.Error(err))
		return &user_service.Administration{}, err
	}

	return resp, nil
}

func (f *UserService) GetList(ctx context.Context, req *user_service.GetListAdministrationRequest) (*user_service.GetListAdministrationResponse, error) {

	f.log.Info("---GetAllAdministration--->>>", logger.Any("req", req))

	resp, err := f.strg.Administration().GetList(ctx, req)
	if err != nil {
		f.log.Error("---GetAllAdministration--->>>", logger.Error(err))
		return &user_service.GetListAdministrationResponse{}, err
	}

	return resp, nil
}

func (f *UserService) Update(ctx context.Context, req *user_service.UpdateAdministration) (*user_service.Administration, error) {
	f.log.Info("---UpdateAdministration--->>>", logger.Any("req", req))

	resp, err := f.strg.Administration().Update(ctx, req)
	if err != nil {
		f.log.Error("---UpdateAdministration--->>>", logger.Error(err))
		return &user_service.Administration{}, err
	}

	return resp, nil
}

func (f *UserService) Delete(ctx context.Context, req *user_service.AdministrationPrimaryKey) (*user_service.EmptyAdmin, error) {
	f.log.Info("---DeleteAdministration--->>>", logger.Any("req", req))

	err := f.strg.Administration().Delete(ctx, req)
	if err != nil {
		f.log.Error("---DeleteAdministration--->>>", logger.Error(err))
		return &user_service.EmptyAdmin{}, err
	}

	return &user_service.EmptyAdmin{}, nil
}

func (s *UserService) GetReportList(ctx context.Context, req *user_service.GetReportListAdministrationRequest) (*user_service.GetReportListAdministrationResponse, error) {
	s.log.Info("---GetReportAdministrationList--->>>", logger.Any("req", req))

	resp, err := s.strg.Administration().GetReportList(ctx, req)
	if err != nil {
		s.log.Error("---GetReportAdministrationList--->>>", logger.Error(err))
		return &user_service.GetReportListAdministrationResponse{}, err
	}

	return resp, nil
}