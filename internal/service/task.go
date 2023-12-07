package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/romandnk/todo/internal/constant"
	"github.com/romandnk/todo/internal/entity"
	storage "github.com/romandnk/todo/internal/repo"
	"github.com/romandnk/todo/pkg/logger"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type taskService struct {
	task   storage.Task
	status storage.Status
	logger logger.Logger
}

func newTaskService(task storage.Task, status storage.Status, logger logger.Logger) *taskService {
	return &taskService{
		task:   task,
		status: status,
		logger: logger,
	}
}

func (s *taskService) CreateTask(ctx context.Context, params CreateTaskParams) (CreateTaskResponse, error) {
	var response CreateTaskResponse

	params.Title = strings.TrimSpace(params.Title)
	params.Description = strings.TrimSpace(params.Description)
	params.StatusName = strings.TrimSpace(params.StatusName)
	params.Date = strings.TrimSpace(params.Date)

	if params.Title == "" {
		return response, constant.ErrEmptyTitle
	}
	if utf8.RuneCountInString(params.Title) > 64 {
		return response, constant.ErrTooLongTitle
	}
	if params.Description == "" {
		return response, constant.ErrEmptyDescription
	}
	if params.StatusName == "" {
		return response, constant.ErrEmptyStatusName
	}
	if params.Date == "" {
		return response, constant.ErrEmptyDate
	}
	date, err := time.Parse(time.RFC3339, params.Date)
	if err != nil {
		s.logger.Error("error parsing date", zap.Error(err))
		return response, constant.ErrInvalidDateFormat
	}
	now := time.Now().UTC()
	if date.UTC().Before(now) {
		return response, constant.ErrOutdatedDate
	}

	status, err := s.status.GetStatusByName(ctx, params.StatusName)
	if err != nil {
		s.logger.Error("error getting repo status by name", zap.Error(err))
		if errors.Is(err, pgx.ErrNoRows) {
			return response, errors.New(fmt.Sprintf("status name '%s' is not found", params.StatusName))
		}
		return response, constant.ErrInternalError
	}

	task := entity.Task{
		Title:       params.Title,
		Description: params.Description,
		StatusID:    status.ID,
		Date:        date.UTC(),
		Deleted:     false,
		CreatedAt:   now,
		DeletedAt:   time.Time{},
	}
	id, err := s.task.CreateTask(ctx, task)
	if err != nil {
		s.logger.Error("error creating repo task", zap.Error(err))
		return response, constant.ErrInternalError
	}

	response.ID = id

	return response, nil
}

func (s *taskService) DeleteTaskByID(ctx context.Context, stringID string) error {
	if stringID == "" {
		return constant.ErrEmptyTaskID
	}
	id, err := strconv.Atoi(stringID)
	if err != nil {
		s.logger.Error("error converting string task id to int task id", zap.Error(err))
		return constant.ErrInvalidTaskID
	}

	err = s.task.DeleteTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, constant.ErrTaskIDNotExists) {
			return err
		}
		s.logger.Error("error deleting repo task by id", zap.Error(err))
		return constant.ErrInternalError
	}

	return nil
}
