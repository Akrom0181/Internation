package service

import (
	"context"
	"schedule_service/config"
	"schedule_service/genproto/schedule_service"
	"schedule_service/grpc/client"
	"schedule_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
)

type StudentPaymentService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
}

func NewStudentPaymentService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *StudentPaymentService {
	return &StudentPaymentService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (s *StudentPaymentService) Create(ctx context.Context, req *schedule_service.CreateStudentPayment) (*schedule_service.GetStudentPayment, error) {
	s.log.Info("---CreateStudentPayment--->>>", logger.Any("req", req))

	resp, err := s.strg.StudentPayment().Create(ctx, req)
	if err != nil {
		s.log.Error("---CreateStudentPayment--->>>", logger.Error(err))
		return &schedule_service.GetStudentPayment{}, err
	}

	return resp, nil
}

func (s *StudentPaymentService) GetByID(ctx context.Context, req *schedule_service.StudentPaymentPrimaryKey) (*schedule_service.GetStudentPayment, error) {
	s.log.Info("---GetSingleStudentPayment--->>>", logger.Any("req", req))

	resp, err := s.strg.StudentPayment().GetByID(ctx, req)
	if err != nil {
		s.log.Error("---GetSingleStudentPayment--->>>", logger.Error(err))
		return &schedule_service.GetStudentPayment{}, err
	}

	return resp, nil
}

func (s *StudentPaymentService) GetList(ctx context.Context, req *schedule_service.GetListStudentPaymentRequest) (*schedule_service.GetListStudentPaymentResponse, error) {
	s.log.Info("---GetAllStudentPayments--->>>", logger.Any("req", req))

	resp, err := s.strg.StudentPayment().GetList(ctx, req)
	if err != nil {
		s.log.Error("---GetAllStudentPayments--->>>", logger.Error(err))
		return &schedule_service.GetListStudentPaymentResponse{}, err
	}

	return resp, nil
}

func (s *StudentPaymentService) Update(ctx context.Context, req *schedule_service.UpdateStudentPayment) (*schedule_service.GetStudentPayment, error) {
	s.log.Info("---UpdateStudentPayment--->>>", logger.Any("req", req))

	resp, err := s.strg.StudentPayment().Update(ctx, req)
	if err != nil {
		s.log.Error("---UpdateStudentPayment--->>>", logger.Error(err))
		return &schedule_service.GetStudentPayment{}, err
	}

	return resp, nil
}

func (s *StudentPaymentService) Delete(ctx context.Context, req *schedule_service.StudentPaymentPrimaryKey) (*schedule_service.EmptyStudentPayment, error) {
	s.log.Info("---DeleteStudentPayment--->>>", logger.Any("req", req))

	err := s.strg.StudentPayment().Delete(ctx, req)
	if err != nil {
		s.log.Error("---DeleteStudentPayment--->>>", logger.Error(err))
		return &schedule_service.EmptyStudentPayment{}, err
	}

	return &schedule_service.EmptyStudentPayment{}, nil
}
