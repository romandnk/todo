package service

//go:generate mockgen -source=service.go -destination=mock/mock.go service

import (
	"context"
	storage "github.com/romandnk/todo/internal/repo"
	statusservice "github.com/romandnk/todo/internal/service/status"
	"github.com/romandnk/todo/internal/service/task"
	"github.com/romandnk/todo/pkg/logger"
)

type Task interface {
	CreateTask(ctx context.Context, params taskservice.CreateTaskParams) (taskservice.CreateTaskResponse, error)
	GetAllTasks(ctx context.Context, limitStr, lastIDStr, statusName, dateStr string) (taskservice.GetAllTasksResponse, error)
	GetTaskByID(ctx context.Context, stringID string) (taskservice.GetTaskWithStatusNameModel, error)
	UpdateTaskByID(ctx context.Context, stringID string, params taskservice.UpdateTaskByIDParams) error
	DeleteTaskByID(ctx context.Context, stringID string) error
}

type Status interface {
	CreateStatus(ctx context.Context, params statusservice.CreateStatusParams) (statusservice.CreateStatusResponse, error)
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
		Status: statusservice.NewStatusService(dep.Repo.Status, dep.Logger),
		Task:   taskservice.NewTaskService(dep.Repo.Task, dep.Repo.Status, dep.Logger),
	}
}
