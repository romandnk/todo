package storage

//go:generate mockgen -source=storage.go -destination=mock/mock.go storage

import (
	"context"
	"github.com/romandnk/todo/internal/entity"
	postgresrepo "github.com/romandnk/todo/internal/repo/postgres"
	postgres "github.com/romandnk/todo/pkg/storage"
)

type Task interface {
	CreateTask(ctx context.Context, task entity.Task) (int, error)
	GetAllTasks(ctx context.Context) ([]*entity.Task, error)
	GetTaskByID(ctx context.Context, id int) (entity.Task, error)
	UpdateTaskByID(ctx context.Context, id int, task entity.Task) error
	DeleteTaskByID(ctx context.Context, id int) error
}

type Status interface {
	CreateStatus(ctx context.Context, status entity.Status) (int, error)
	GetAllStatuses(ctx context.Context) ([]*entity.Status, error)
	GetStatusByName(ctx context.Context, name string) (entity.Status, error)
	GetStatusByID(ctx context.Context, id int) (entity.Status, error)
}

type Repository struct {
	Task   Task
	Status Status
}

func NewRepository(db postgres.PgxPool) *Repository {
	return &Repository{
		Task:   postgresrepo.NewTaskRepo(db),
		Status: postgresrepo.NewStatusRepo(db),
	}
}
