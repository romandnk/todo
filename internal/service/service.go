package service

import (
	"context"
	storage "github.com/romandnk/todo/internal/repo"
	"github.com/romandnk/todo/pkg/logger"
)

type Task interface {
	CreateTask(ctx context.Context, params CreateTaskParams) (CreateTaskResponse, error)
	GetTaskByID(ctx context.Context, stringID string) (GetTaskWithStatusNameModel, error)
	UpdateTaskByID(ctx context.Context, stringID string, params UpdateTaskByIDParams) error
	DeleteTaskByID(ctx context.Context, stringID string) error
}

type Status interface {
	CreateStatus(ctx context.Context, params CreateStatusParams) (CreateStatusResponse, error)
}

type Services struct {
	Status Status
	Task   Task
}

type Dependencies struct {
	Repo   *storage.Repository
	Logger logger.Logger
}

func NewServices(dep Dependencies) *Services {
	return &Services{
		Status: newStatusService(dep.Repo.Status, dep.Logger),
		Task:   newTaskService(dep.Repo.Task, dep.Repo.Status, dep.Logger),
	}
}
