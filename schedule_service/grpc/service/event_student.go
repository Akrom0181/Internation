package service

import (
	"context"
	"schedule_service/config"
	"schedule_service/genproto/schedule_service"
	"schedule_service/grpc/client"
	"schedule_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
)

type EventStudentService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
}

func NewEventStudentService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *EventStudentService {
	return &EventStudentService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (f *EventStudentService) Create(ctx context.Context, req *schedule_service.CreateEventStudent) (*schedule_service.GetEventStudent, error) {

	f.log.Info("---CreateEventStudent--->>>", logger.Any("req", req))

	resp, err := f.strg.EventStudent().Create(ctx, req)
	if err != nil {
		f.log.Error("---CreateEventStudent--->>>", logger.Error(err))
		return &schedule_service.GetEventStudent{}, err
	}

	return resp, nil
}

func (f *EventStudentService) GetByID(ctx context.Context, req *schedule_service.EventStudentPrimaryKey) (*schedule_service.GetEventStudent, error) {
	f.log.Info("---GetSingleEveCreateEventStudent--->>>", logger.Any("req", req))

	resp, err := f.strg.EventStudent().GetByID(ctx, req)
	if err != nil {
		f.log.Error("---GetSingleEveCreateEventStudent--->>>", logger.Error(err))
		return &schedule_service.GetEventStudent{}, err
	}

	return resp, nil
}

func (f *EventStudentService) GetList(ctx context.Context, req *schedule_service.GetListEventStudentRequest) (*schedule_service.GetListEventStudentResponse, error) {

	f.log.Info("---GetAllEventStudent--->>>", logger.Any("req", req))

	resp, err := f.strg.EventStudent().GetList(ctx, req)
	if err != nil {
		f.log.Error("---GetAllEventStudent--->>>", logger.Error(err))
		return &schedule_service.GetListEventStudentResponse{}, err
	}

	return resp, nil
}

func (f *EventStudentService) Update(ctx context.Context, req *schedule_service.UpdateEventStudent) (*schedule_service.GetEventStudent, error) {
	f.log.Info("---UpdateEveCreateEventStudent--->>>", logger.Any("req", req))

	resp, err := f.strg.EventStudent().Update(ctx, req)
	if err != nil {
		f.log.Error("---UpdateEveCreateEventStudent--->>>", logger.Error(err))
		return &schedule_service.GetEventStudent{}, err
	}

	return resp, nil
}

func (f *EventStudentService) Delete(ctx context.Context, req *schedule_service.EventStudentPrimaryKey) (*schedule_service.EmptyEventStudent, error) {
	f.log.Info("---DeleteEveCreateEventStudent--->>>", logger.Any("req", req))

	err := f.strg.EventStudent().Delete(ctx, req)
	if err != nil {
		f.log.Error("---DeleteEveCreateEventStudent--->>>", logger.Error(err))
		return &schedule_service.EmptyEventStudent{}, err
	}

	return &schedule_service.EmptyEventStudent{}, nil
}

func (s *EventStudentService) GetStudentByID(ctx context.Context, req *schedule_service.EventStudentPrimaryKey) (*schedule_service.GetStudentWithEventsResponse, error) {
	s.log.Info("---GetStudentByID--->>>", logger.Any("req", req))

	resp, err := s.strg.EventStudent().GetStudentWithEventsByID(ctx, req)
	if err != nil {
		s.log.Error("---GetStudentByID--->>>", logger.Error(err))
		return &schedule_service.GetStudentWithEventsResponse{}, err
	}
	return resp, nil
}