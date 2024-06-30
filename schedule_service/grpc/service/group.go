package service

import (
	"context"
	"schedule_service/config"
	"schedule_service/genproto/schedule_service"
	"schedule_service/grpc/client"
	"schedule_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
)

type GroupService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
}

func NewGroupService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *GroupService {
	return &GroupService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (f *GroupService) Create(ctx context.Context, req *schedule_service.CreateGroup) (*schedule_service.GetGroup, error) {

	f.log.Info("---CreateGroup--->>>", logger.Any("req", req))

	resp, err := f.strg.Group().Create(ctx, req)
	if err != nil {
		f.log.Error("---CreateGroup--->>>", logger.Error(err))
		return &schedule_service.GetGroup{}, err
	}

	return resp, nil
}

func (f *GroupService) GetByID(ctx context.Context, req *schedule_service.GroupPrimaryKey) (*schedule_service.GetGroup, error) {
	f.log.Info("---GetSingleGroup--->>>", logger.Any("req", req))

	resp, err := f.strg.Group().GetByID(ctx, req)
	if err != nil {
		f.log.Error("---GetSingleGroup--->>>", logger.Error(err))
		return &schedule_service.GetGroup{}, err
	}

	return resp, nil
}

func (f *GroupService) GetList(ctx context.Context, req *schedule_service.GetListGroupRequest) (*schedule_service.GetListGroupResponse, error) {

	f.log.Info("---GetAllGroup--->>>", logger.Any("req", req))

	resp, err := f.strg.Group().GetList(ctx, req)
	if err != nil {
		f.log.Error("---GetAllGroup--->>>", logger.Error(err))
		return &schedule_service.GetListGroupResponse{}, err
	}

	return resp, nil
}

func (f *GroupService) Update(ctx context.Context, req *schedule_service.UpdateGroup) (*schedule_service.GetGroup, error) {
	f.log.Info("---UpdateGroup--->>>", logger.Any("req", req))

	resp, err := f.strg.Group().Update(ctx, req)
	if err != nil {
		f.log.Error("---UpdateGroup--->>>", logger.Error(err))
		return &schedule_service.GetGroup{}, err
	}

	return resp, nil
}

func (f *GroupService) Delete(ctx context.Context, req *schedule_service.GroupPrimaryKey) (*schedule_service.EmptyGroup, error) {
	f.log.Info("---DeleteGroup--->>>", logger.Any("req", req))

	err := f.strg.Group().Delete(ctx, req)
	if err != nil {
		f.log.Error("---DeleteGroup--->>>", logger.Error(err))
		return &schedule_service.EmptyGroup{}, err
	}

	return &schedule_service.EmptyGroup{}, nil
}

func (s *GroupService) GetByIDTeacher(ctx context.Context, req *schedule_service.TeacherID) (*schedule_service.GetGroup, error) {
	s.log.Info("---GetGroupByIDTeacher--->>>", logger.Any("req", req))

	resp, err := s.strg.Group().GetByIDTeacher(ctx, req)
	if err != nil {
		s.log.Error("---GetGroupByIDTeacher--->>>", logger.Error(err))
		return &schedule_service.GetGroup{}, err
	}

	return resp, nil
}