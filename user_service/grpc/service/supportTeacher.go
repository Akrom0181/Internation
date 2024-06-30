package service

import (
	"context"
	"user_service/config"
	"user_service/genproto/user_service"
	"user_service/grpc/client"

	"user_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
)

type SupportTeacherService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
}

func NewSupportTeacherService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *SupportTeacherService {
	return &SupportTeacherService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (s *SupportTeacherService) Create(ctx context.Context, req *user_service.CreateSupportTeacher) (*user_service.SupportTeacher, error) {
	s.log.Info("---CreateSupportTeacher--->>>", logger.Any("req", req))

	resp, err := s.strg.SupportTeacher().Create(ctx, req)
	if err != nil {
		s.log.Error("---CreateSupportTeacher--->>>", logger.Error(err))
		return &user_service.SupportTeacher{}, err
	}

	return resp, nil
}

func (s *SupportTeacherService) GetByID(ctx context.Context, req *user_service.SupportTeacherPrimaryKey) (*user_service.SupportTeacher, error) {
	s.log.Info("---GetSingleSupportTeacher--->>>", logger.Any("req", req))

	resp, err := s.strg.SupportTeacher().GetByID(ctx, req)
	if err != nil {
		s.log.Error("---GetSingleSupportTeacher--->>>", logger.Error(err))
		return &user_service.SupportTeacher{}, err
	}

	return resp, nil
}

func (s *SupportTeacherService) GetList(ctx context.Context, req *user_service.GetListSupportTeacherRequest) (*user_service.GetListSupportTeacherResponse, error) {
	s.log.Info("---GetAllSupportTeachers--->>>", logger.Any("req", req))

	resp, err := s.strg.SupportTeacher().GetList(ctx, req)
	if err != nil {
		s.log.Error("---GetAllSupportTeachers--->>>", logger.Error(err))
		return &user_service.GetListSupportTeacherResponse{}, err
	}

	return resp, nil
}

func (s *SupportTeacherService) Update(ctx context.Context, req *user_service.UpdateSupportTeacher) (*user_service.SupportTeacher, error) {
	s.log.Info("---UpdateSupportTeacher--->>>", logger.Any("req", req))

	resp, err := s.strg.SupportTeacher().Update(ctx, req)
	if err != nil {
		s.log.Error("---UpdateSupportTeacher--->>>", logger.Error(err))
		return &user_service.SupportTeacher{}, err
	}

	return resp, nil
}

func (s *SupportTeacherService) Delete(ctx context.Context, req *user_service.SupportTeacherPrimaryKey) (*user_service.EmptySTeacher, error) {
	s.log.Info("---DeleteSupportTeacher--->>>", logger.Any("req", req))

	err := s.strg.SupportTeacher().Delete(ctx, req)
	if err != nil {
		s.log.Error("---DeleteSupportTeacher--->>>", logger.Error(err))
		return &user_service.EmptySTeacher{}, err
	}

	return &user_service.EmptySTeacher{}, nil
}

func (s *SupportTeacherService) GetReportList(ctx context.Context, req *user_service.GetReportListSupportTeacherRequest) (*user_service.GetReportListSupportTeacherResponse, error) {
	s.log.Info("---GetSupportTeacherList--->>>", logger.Any("req", req))

	resp, err := s.strg.SupportTeacher().GetReportList(ctx, req)
	if err != nil {
		s.log.Error("---GetSupportTeacherList--->>>", logger.Error(err))
		return &user_service.GetReportListSupportTeacherResponse{}, err
	}

	return resp, nil
}