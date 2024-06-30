package service

import (
	"context"
	"schedule_service/config"
	"schedule_service/genproto/schedule_service"
	"schedule_service/grpc/client"
	"schedule_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
)

type StudentTaskService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
}

func NewStudentTaskService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *StudentTaskService {
	return &StudentTaskService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (s *StudentTaskService) Create(ctx context.Context, req *schedule_service.CreateStudentTask) (*schedule_service.GetStudentTask, error) {
	s.log.Info("---CreateStudentTask--->>>", logger.Any("req", req))

	resp, err := s.strg.StudentTask().Create(ctx, req)
	if err != nil {
		s.log.Error("---CreateStudentTask--->>>", logger.Error(err))
		return &schedule_service.GetStudentTask{}, err
	}

	return resp, nil
}

func (s *StudentTaskService) GetByID(ctx context.Context, req *schedule_service.StudentTaskPrimaryKey) (*schedule_service.GetStudentTask, error) {
	s.log.Info("---GetSingleStudentTask--->>>", logger.Any("req", req))

	resp, err := s.strg.StudentTask().GetByID(ctx, req)
	if err != nil {
		s.log.Error("---GetSingleStudentTask--->>>", logger.Error(err))
		return &schedule_service.GetStudentTask{}, err
	}

	return resp, nil
}

func (s *StudentTaskService) GetList(ctx context.Context, req *schedule_service.GetListStudentTaskRequest) (*schedule_service.GetListStudentTaskResponse, error) {
	s.log.Info("---GetAllStudentTasks--->>>", logger.Any("req", req))

	resp, err := s.strg.StudentTask().GetList(ctx, req)
	if err != nil {
		s.log.Error("---GetAllStudentTasks--->>>", logger.Error(err))
		return &schedule_service.GetListStudentTaskResponse{}, err
	}

	return resp, nil
}

func (s *StudentTaskService) Update(ctx context.Context, req *schedule_service.UpdateStudentTask) (*schedule_service.GetStudentTask, error) {
	s.log.Info("---UpdateStudentTask--->>>", logger.Any("req", req))

	resp, err := s.strg.StudentTask().Update(ctx, req)
	if err != nil {
		s.log.Error("---UpdateStudentTask--->>>", logger.Error(err))
		return &schedule_service.GetStudentTask{}, err
	}

	return resp, nil
}

func (s *StudentTaskService) Delete(ctx context.Context, req *schedule_service.StudentTaskPrimaryKey) (*schedule_service.EmptyStudentTask, error) {
	s.log.Info("---DeleteStudentTask--->>>", logger.Any("req", req))

	err := s.strg.StudentTask().Delete(ctx, req)
	if err != nil {
		s.log.Error("---DeleteStudentTask--->>>", logger.Error(err))
		return &schedule_service.EmptyStudentTask{}, err
	}

	return &schedule_service.EmptyStudentTask{}, nil
}

func (s *StudentTaskService) UpdateScoreforTeacher(ctx context.Context, req *schedule_service.UpdateStudentScoreRequest) (*schedule_service.GetStudentTask, error) {
	s.log.Info("---UpdateStudentTask--->>>", logger.Any("req", req))
	resp, err := s.strg.StudentTask().UpdateScoreforTeacher(ctx, req)
	if err != nil {
		s.log.Error("---UpdateStudentTask--->>>", logger.Error(err))
		return &schedule_service.GetStudentTask{}, err
	}
	return resp, nil
}

func (s *StudentTaskService)UpdateScoreforStudent(ctx context.Context, req *schedule_service.UpdateStudentScoreRequest) (*schedule_service.GetStudentTask, error) {
	s.log.Info("---UpdateStudentTask--->>>", logger.Any("req", req))
	resp, err := s.strg.StudentTask().UpdateScoreforStudent(ctx, req)
	if err != nil {
		s.log.Error("---UpdateStudentTask--->>>", logger.Error(err))
		return &schedule_service.GetStudentTask{}, err
	}
	return resp, nil
}