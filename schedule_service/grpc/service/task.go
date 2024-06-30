package service

import (
	"context"
	"schedule_service/config"
	"schedule_service/genproto/schedule_service"
	"schedule_service/grpc/client"
	"schedule_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
)

type TaskService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
}

func NewTaskService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *TaskService {
	return &TaskService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (t *TaskService) Create(ctx context.Context, req *schedule_service.CreateTask) (*schedule_service.GetTask, error) {
	t.log.Info("---CreateTask--->>>", logger.Any("req", req))

	resp, err := t.strg.Task().Create(ctx, req)
	if err != nil {
		t.log.Error("---CreateTask--->>>", logger.Error(err))
		return &schedule_service.GetTask{}, err
	}

	return resp, nil
}

func (t *TaskService) GetByID(ctx context.Context, req *schedule_service.TaskPrimaryKey) (*schedule_service.GetTask, error) {
	t.log.Info("---GetSingleTask--->>>", logger.Any("req", req))

	resp, err := t.strg.Task().GetByID(ctx, req)
	if err != nil {
		t.log.Error("---GetSingleTask--->>>", logger.Error(err))
		return &schedule_service.GetTask{}, err
	}

	return resp, nil
}

func (t *TaskService) GetList(ctx context.Context, req *schedule_service.GetListTaskRequest) (*schedule_service.GetListTaskResponse, error) {
	t.log.Info("---GetAllTasks--->>>", logger.Any("req", req))

	resp, err := t.strg.Task().GetList(ctx, req)
	if err != nil {
		t.log.Error("---GetAllTasks--->>>", logger.Error(err))
		return &schedule_service.GetListTaskResponse{}, err
	}

	return resp, nil
}

func (t *TaskService) Update(ctx context.Context, req *schedule_service.UpdateTask) (*schedule_service.GetTask, error) {
	t.log.Info("---UpdateTask--->>>", logger.Any("req", req))

	resp, err := t.strg.Task().Update(ctx, req)
	if err != nil {
		t.log.Error("---UpdateTask--->>>", logger.Error(err))
		return &schedule_service.GetTask{}, err
	}

	return resp, nil
}

func (t *TaskService) Delete(ctx context.Context, req *schedule_service.TaskPrimaryKey) (*schedule_service.EmptyTask, error) {
	t.log.Info("---DeleteTask--->>>", logger.Any("req", req))

	err := t.strg.Task().Delete(ctx, req)
	if err != nil {
		t.log.Error("---DeleteTask--->>>", logger.Error(err))
		return &schedule_service.EmptyTask{}, err
	}

	return &schedule_service.EmptyTask{}, nil
}
