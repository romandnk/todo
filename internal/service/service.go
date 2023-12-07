package service

import (
	"context"
	storage "github.com/romandnk/todo/internal/repo"
	"github.com/romandnk/todo/pkg/logger"
)

type Task interface {
}

type Status interface {
	CreateStatus(ctx context.Context, params CreateStatusParams) (CreateStatusResponse, error)
}

type Services struct {
	Status Status
}

type Dependencies struct {
	Repo   *storage.Repository
	Logger logger.Logger
}

func NewServices(dep Dependencies) *Services {
	return &Services{
		Status: newStatusService(dep.Repo.Status, dep.Logger),
	}
}
