package service

import (
	"context"
	"schedule_service/config"
	"schedule_service/genproto/schedule_service"
	"schedule_service/grpc/client"
	"schedule_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
)

type EventService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
}

func NewEventService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *EventService {
	return &EventService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (f *EventService) Create(ctx context.Context, req *schedule_service.CreateEvent) (*schedule_service.GetEvent, error) {

	f.log.Info("---CreateEvent--->>>", logger.Any("req", req))

	resp, err := f.strg.Event().Create(ctx, req)
	if err != nil {
		f.log.Error("---CreateEvent--->>>", logger.Error(err))
		return &schedule_service.GetEvent{}, err
	}

	return resp, nil
}

func (f *EventService) GetByID(ctx context.Context, req *schedule_service.EventPrimaryKey) (*schedule_service.GetEvent, error) {
	f.log.Info("---GetSingleEvent--->>>", logger.Any("req", req))

	resp, err := f.strg.Event().GetByID(ctx, req)
	if err != nil {
		f.log.Error("---GetSingleEvent--->>>", logger.Error(err))
		return &schedule_service.GetEvent{}, err
	}

	return resp, nil
}

func (f *EventService) GetList(ctx context.Context, req *schedule_service.GetListEventRequest) (*schedule_service.GetListEventResponse, error) {

	f.log.Info("---GetAllEvent--->>>", logger.Any("req", req))

	resp, err := f.strg.Event().GetList(ctx, req)
	if err != nil {
		f.log.Error("---GetAllEvent--->>>", logger.Error(err))
		return &schedule_service.GetListEventResponse{}, err
	}

	return resp, nil
}

func (f *EventService) Update(ctx context.Context, req *schedule_service.UpdateEvent) (*schedule_service.GetEvent, error) {
	f.log.Info("---UpdateEvent--->>>", logger.Any("req", req))

	resp, err := f.strg.Event().Update(ctx, req)
	if err != nil {
		f.log.Error("---UpdateEvent--->>>", logger.Error(err))
		return &schedule_service.GetEvent{}, err
	}

	return resp, nil
}

func (f *EventService) Delete(ctx context.Context, req *schedule_service.EventPrimaryKey) (*schedule_service.EmptyEvent, error) {
	f.log.Info("---DeleteEvent--->>>", logger.Any("req", req))

	err := f.strg.Event().Delete(ctx, req)
	if err != nil {
		f.log.Error("---DeleteEvent--->>>", logger.Error(err))
		return &schedule_service.EmptyEvent{}, err
	}

	return &schedule_service.EmptyEvent{}, nil
}
