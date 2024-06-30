package service

import (
	"context"
	"schedule_service/config"
	"schedule_service/genproto/schedule_service"
	"schedule_service/grpc/client"
	"schedule_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
)

type ScheduleService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
}

func NewScheduleService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *ScheduleService {
	return &ScheduleService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (s *ScheduleService) Create(ctx context.Context, req *schedule_service.CreateSchedule) (*schedule_service.GetSchedule, error) {
	s.log.Info("---CreateSchedule--->>>", logger.Any("req", req))

	resp, err := s.strg.Schedule().Create(ctx, req)
	if err != nil {
		s.log.Error("---CreateSchedule--->>>", logger.Error(err))
		return &schedule_service.GetSchedule{}, err
	}

	return resp, nil
}

func (s *ScheduleService) GetByID(ctx context.Context, req *schedule_service.SchedulePrimaryKey) (*schedule_service.GetSchedule, error) {
	s.log.Info("---GetSingleSchedule--->>>", logger.Any("req", req))

	resp, err := s.strg.Schedule().GetByID(ctx, req)
	if err != nil {
		s.log.Error("---GetSingleSchedule--->>>", logger.Error(err))
		return &schedule_service.GetSchedule{}, err
	}

	return resp, nil
}

func (s *ScheduleService) GetList(ctx context.Context, req *schedule_service.GetListScheduleRequest) (*schedule_service.GetListScheduleResponse, error) {
	s.log.Info("---GetAllSchedules--->>>", logger.Any("req", req))

	resp, err := s.strg.Schedule().GetList(ctx, req)
	if err != nil {
		s.log.Error("---GetAllSchedules--->>>", logger.Error(err))
		return &schedule_service.GetListScheduleResponse{}, err
	}

	return resp, nil
}

func (s *ScheduleService) Update(ctx context.Context, req *schedule_service.UpdateSchedule) (*schedule_service.GetSchedule, error) {
	s.log.Info("---UpdateSchedule--->>>", logger.Any("req", req))

	resp, err := s.strg.Schedule().Update(ctx, req)
	if err != nil {
		s.log.Error("---UpdateSchedule--->>>", logger.Error(err))
		return &schedule_service.GetSchedule{}, err
	}

	return resp, nil
}

func (s *ScheduleService) Delete(ctx context.Context, req *schedule_service.SchedulePrimaryKey) (*schedule_service.EmptySchedule, error) {
	s.log.Info("---DeleteSchedule--->>>", logger.Any("req", req))

	err := s.strg.Schedule().Delete(ctx, req)
	if err != nil {
		s.log.Error("---DeleteSchedule--->>>", logger.Error(err))
		return &schedule_service.EmptySchedule{}, err
	}

	return &schedule_service.EmptySchedule{}, nil
}

func (s *ScheduleService) GetScheduleForMonth(ctx context.Context, req *schedule_service.GetScheduleForMonthRequest) (*schedule_service.GetListScheduleResponse, error) {
    s.log.Info("---GetScheduleForMonth--->>>", logger.Any("req", req))

    resp, err := s.strg.Schedule().GetScheduleForMonth(ctx, req.TeacherId, req.MonthStartDate, req.MonthEndDate)
    if err != nil {
        s.log.Error("---GetScheduleForMonth--->>>", logger.Error(err))
        return &schedule_service.GetListScheduleResponse{}, err
    }

    return resp, nil
}

func (s *ScheduleService) GetScheduleForWeek(ctx context.Context, req *schedule_service.GetScheduleForWeekRequest) (*schedule_service.GetListScheduleResponse, error) {
    s.log.Info("---GetScheduleForWeek--->>>", logger.Any("req", req))

    resp, err := s.strg.Schedule().GetScheduleForMonth(ctx, req.TeacherId, req.WeekStartDate, req.WeekEndDate)
    if err != nil {
        s.log.Error("---GetScheduleForWeek--->>>", logger.Error(err))
        return &schedule_service.GetListScheduleResponse{}, err
    }

    return resp, nil
}