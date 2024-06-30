package service

import (
	"context"
	"schedule_service/config"
	"schedule_service/genproto/schedule_service"
	"schedule_service/grpc/client"
	"schedule_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
)

type JournalService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
}

func NewJournalService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *JournalService {
	return &JournalService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (j *JournalService) Create(ctx context.Context, req *schedule_service.CreateJournal) (*schedule_service.GetJournal, error) {
	j.log.Info("---CreateJournal--->>>", logger.Any("req", req))

	resp, err := j.strg.Journal().Create(ctx, req)
	if err != nil {
		j.log.Error("---CreateJournal--->>>", logger.Error(err))
		return &schedule_service.GetJournal{}, err
	}

	return resp, nil
}

func (j *JournalService) GetByID(ctx context.Context, req *schedule_service.JournalPrimaryKey) (*schedule_service.GetJournal, error) {
	j.log.Info("---GetSingleJournal--->>>", logger.Any("req", req))

	resp, err := j.strg.Journal().GetByID(ctx, req)
	if err != nil {
		j.log.Error("---GetSingleJournal--->>>", logger.Error(err))
		return &schedule_service.GetJournal{}, err
	}

	return resp, nil
}

func (j *JournalService) GetList(ctx context.Context, req *schedule_service.GetListJournalRequest) (*schedule_service.GetListJournalResponse, error) {
	j.log.Info("---GetAllJournals--->>>", logger.Any("req", req))

	resp, err := j.strg.Journal().GetList(ctx, req)
	if err != nil {
		j.log.Error("---GetAllJournals--->>>", logger.Error(err))
		return &schedule_service.GetListJournalResponse{}, err
	}

	return resp, nil
}

func (j *JournalService) Update(ctx context.Context, req *schedule_service.UpdateJournal) (*schedule_service.GetJournal, error) {
	j.log.Info("---UpdateJournal--->>>", logger.Any("req", req))

	resp, err := j.strg.Journal().Update(ctx, req)
	if err != nil {
		j.log.Error("---UpdateJournal--->>>", logger.Error(err))
		return &schedule_service.GetJournal{}, err
	}

	return resp, nil
}

func (j *JournalService) Delete(ctx context.Context, req *schedule_service.JournalPrimaryKey) (*schedule_service.EmptyJournal, error) {
	j.log.Info("---DeleteJournal--->>>", logger.Any("req", req))

	err := j.strg.Journal().Delete(ctx, req)
	if err != nil {
		j.log.Error("---DeleteJournal--->>>", logger.Error(err))
		return &schedule_service.EmptyJournal{}, err
	}

	return &schedule_service.EmptyJournal{}, nil
}

func (s *JournalService) GetByIDStudent(ctx context.Context, req *schedule_service.JournalPrimaryKey) (*schedule_service.GetJournal, error) {
	s.log.Info("---GetJurnalByID--->>>", logger.Any("req", req))

	resp, err := s.strg.Journal().GetByGroupID(ctx, req)
	if err != nil {
		s.log.Error("---GetJurnalByID--->>>", logger.Error(err))
		return &schedule_service.GetJournal{}, err
	}

	return resp, nil
}