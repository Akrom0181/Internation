package service

import (
	"context"
	"user_service/config"
	"user_service/genproto/user_service"
	"user_service/grpc/client"

	"user_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
)

type StudentService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
}

func NewStudentService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *StudentService {
	return &StudentService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (s *StudentService) Create(ctx context.Context, req *user_service.CreateStudent) (*user_service.Student, error) {
	s.log.Info("---CreateStudent--->>>", logger.Any("req", req))

	resp, err := s.strg.Student().Create(ctx, req)
	if err != nil {
		s.log.Error("---CreateStudent--->>>", logger.Error(err))
		return &user_service.Student{}, err
	}

	return resp, nil
}

func (s *StudentService) GetByID(ctx context.Context, req *user_service.StudentPrimaryKey) (*user_service.Student, error) {
	s.log.Info("---GetSingleStudent--->>>", logger.Any("req", req))

	resp, err := s.strg.Student().GetByID(ctx, req)
	if err != nil {
		s.log.Error("---GetSingleStudent--->>>", logger.Error(err))
		return &user_service.Student{}, err
	}

	return resp, nil
}

func (s *StudentService) GetList(ctx context.Context, req *user_service.GetListStudentRequest) (*user_service.GetListStudentResponse, error) {
	s.log.Info("---GetAllStudents--->>>", logger.Any("req", req))

	resp, err := s.strg.Student().GetList(ctx, req)
	if err != nil {
		s.log.Error("---GetAllStudents--->>>", logger.Error(err))
		return &user_service.GetListStudentResponse{}, err
	}

	return resp, nil
}

func (s *StudentService) Update(ctx context.Context, req *user_service.UpdateStudent) (*user_service.Student, error) {
	s.log.Info("---UpdateStudent--->>>", logger.Any("req", req))

	resp, err := s.strg.Student().Update(ctx, req)
	if err != nil {
		s.log.Error("---UpdateStudent--->>>", logger.Error(err))
		return &user_service.Student{}, err
	}

	return resp, nil
}

func (s *StudentService) Delete(ctx context.Context, req *user_service.StudentPrimaryKey) (*user_service.StudentEmpty, error) {
	s.log.Info("---DeleteStudent--->>>", logger.Any("req", req))

	err := s.strg.Student().Delete(ctx, req)
	if err != nil {
		s.log.Error("---DeleteStudent--->>>", logger.Error(err))
		return &user_service.StudentEmpty{}, err
	}

	return &user_service.StudentEmpty{}, nil
}

func (s *StudentService) GetReportList(ctx context.Context, req *user_service.GetReportListStudentRequest) (*user_service.GetReportListStudentResponse, error) {
	s.log.Info("---GetStudentList--->>>", logger.Any("req", req))

	resp, err := s.strg.Student().GetReportList(ctx, req)
	if err != nil {
		s.log.Error("---GetStudentList--->>>", logger.Error(err))
		return &user_service.GetReportListStudentResponse{}, err
	}

	return resp, nil
}