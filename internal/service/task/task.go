package taskservice

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

type TaskService struct {
	task   storage.Task
	status storage.Status
	logger logger.Logger
}

func NewTaskService(task storage.Task, status storage.Status, logger logger.Logger) *TaskService {
	return &TaskService{
		task:   task,
		status: status,
		logger: logger,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, params CreateTaskParams) (CreateTaskResponse, error) {
	var response CreateTaskResponse

	params.Title = strings.TrimSpace(params.Title)
	params.Description = strings.TrimSpace(params.Description)
	params.StatusName = strings.ToLower(strings.TrimSpace(params.StatusName))
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

func (s *TaskService) DeleteTaskByID(ctx context.Context, stringID string) error {
	if stringID == "" {
		return constant.ErrEmptyTaskID
	}
	id, err := strconv.Atoi(stringID)
	if err != nil {
		s.logger.Error("error converting string task id to int task id", zap.Error(err))
		return constant.ErrInvalidTaskID
	}

	if id <= 0 {
		return constant.ErrNonPositiveTaskID
	}

	err = s.task.DeleteTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, constant.ErrTaskIDNotExists) {
			return errors.New(fmt.Sprintf("%s %d", err.Error(), id))
		}
		s.logger.Error("error deleting repo task by id", zap.Error(err))
		return constant.ErrInternalError
	}

	return nil
}

func (s *TaskService) UpdateTaskByID(ctx context.Context, stringID string, params UpdateTaskByIDParams) error {
	if stringID == "" {
		return constant.ErrEmptyTaskID
	}
	id, err := strconv.Atoi(stringID)
	if err != nil {
		s.logger.Error("error converting string task id to int task id", zap.Error(err))
		return constant.ErrInvalidTaskID
	}

	if id <= 0 {
		return constant.ErrNonPositiveTaskID
	}

	params.Title = strings.TrimSpace(params.Title)
	params.Description = strings.TrimSpace(params.Description)
	params.StatusName = strings.ToLower(strings.TrimSpace(params.StatusName))
	var status entity.Status
	if params.StatusName != "" {
		status, err = s.status.GetStatusByName(ctx, params.StatusName)
		if err != nil {
			s.logger.Error("error getting repo status by name", zap.Error(err))
			if errors.Is(err, pgx.ErrNoRows) {
				return errors.New(fmt.Sprintf("status name '%s' is not found", params.StatusName))
			}
			return constant.ErrInternalError
		}
	}

	var date time.Time
	params.Date = strings.TrimSpace(params.Date)
	if params.Date != "" {
		date, err = time.Parse(time.RFC3339, params.Date)
		if err != nil {
			s.logger.Error("error parsing date", zap.Error(err))
			return constant.ErrInvalidDateFormat
		}
		now := time.Now().UTC()
		if date.UTC().Before(now) {
			return constant.ErrOutdatedDate
		}
	}

	task := entity.Task{
		Title:       params.Title,
		Description: params.Description,
		StatusID:    status.ID,
		Date:        date,
	}
	err = s.task.UpdateTaskByID(ctx, id, task)
	if err != nil {
		if errors.Is(err, constant.ErrTaskIDNotExists) {
			return errors.New(fmt.Sprintf("%s %d", err.Error(), id))
		}
		s.logger.Error("error updating repo task by id", zap.Error(err))
		return constant.ErrInternalError
	}

	return nil
}

func (s *TaskService) GetAllTasks(ctx context.Context, limitStr, lastIDStr, statusName, dateStr string) (GetAllTasksResponse, error) {
	var response GetAllTasksResponse

	var limit int
	var err error
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			s.logger.Error("error converting limit into int", zap.Error(err))
			return response, constant.ErrInvalidLimit
		}
	}
	if limit < 0 {
		return response, constant.ErrNegativeLimit
	}

	var lastID int
	if lastIDStr != "" {
		lastID, err = strconv.Atoi(lastIDStr)
		if err != nil {
			s.logger.Error("error converting last id into int", zap.Error(err))
			return response, constant.ErrInvalidLastTaskID
		}
	}
	if lastID < 0 {
		return response, constant.ErrNegativeLastTaskID
	}

	var status entity.Status
	var mapStatuses = make(map[int]string)
	statusName = strings.ToLower(statusName)
	if statusName != "" {
		status, err = s.status.GetStatusByName(ctx, statusName)
		if err != nil {
			s.logger.Error("error getting repo status by name", zap.Error(err))
			if errors.Is(err, pgx.ErrNoRows) {
				return response, errors.New(fmt.Sprintf("status name '%s' is not found", statusName))
			}
			return response, constant.ErrInternalError
		}
	} else {
		statuses, err := s.status.GetAllStatuses(ctx)
		if err != nil {
			s.logger.Error("error getting repo all statuses", zap.Error(err))
			if errors.Is(err, pgx.ErrNoRows) {
				return response, errors.New("statuses are not found")
			}
			return response, constant.ErrInternalError
		}

		mapStatuses = make(map[int]string, len(statuses))
		for _, status := range statuses {
			if status != nil {
				mapStatuses[status.ID] = status.Name
			} else {
				s.logger.Error("error status is nil")
				return response, constant.ErrInternalError
			}
		}
	}

	var date time.Time
	if dateStr != "" {
		date, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			s.logger.Error("error parsing date", zap.Error(err))
			return response, constant.ErrInvalidDateFormat
		}
	}

	date = date.UTC()

	tasks, err := s.task.GetAllTasks(ctx, status.ID, limit, lastID, date)
	if err != nil {
		s.logger.Error("error getting repo all tasks", zap.Error(err))
		if errors.Is(err, pgx.ErrNoRows) {
			return response, nil
		}
		return response, constant.ErrInternalError
	}

	response.Tasks = make([]GetTaskWithStatusNameModel, 0, len(tasks))
	taskWithStatusName := GetTaskWithStatusNameModel{}
	for _, task := range tasks {
		if status.ID == 0 {
			taskWithStatusName = GetTaskWithStatusNameModel{
				ID:          task.ID,
				Title:       task.Title,
				Description: task.Description,
				StatusName:  mapStatuses[task.StatusID],
				Date:        task.Date.Format(time.RFC3339),
				CreatedAt:   task.CreatedAt.Format(time.RFC3339),
			}
		} else {
			taskWithStatusName = GetTaskWithStatusNameModel{
				ID:          task.ID,
				Title:       task.Title,
				Description: task.Description,
				StatusName:  status.Name,
				Date:        task.Date.Format(time.RFC3339),
				CreatedAt:   task.CreatedAt.Format(time.RFC3339),
			}
		}

		response.Tasks = append(response.Tasks, taskWithStatusName)
	}

	response.Total = len(response.Tasks)

	return response, nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, stringID string) (GetTaskWithStatusNameModel, error) {
	var response GetTaskWithStatusNameModel

	if stringID == "" {
		return response, constant.ErrEmptyTaskID
	}
	id, err := strconv.Atoi(stringID)
	if err != nil {
		s.logger.Error("error converting string task id to int task id", zap.Error(err))
		return response, constant.ErrInvalidTaskID
	}

	if id <= 0 {
		return response, constant.ErrNonPositiveTaskID
	}

	task, err := s.task.GetTaskByID(ctx, id)
	if err != nil {
		s.logger.Error("error getting repo task by id", zap.Error(err))
		if errors.Is(err, pgx.ErrNoRows) {
			return response, errors.New(fmt.Sprintf("task with id '%d' is not found", id))
		}
		return response, constant.ErrInternalError
	}

	status, err := s.status.GetStatusByID(ctx, task.StatusID)
	if err != nil {
		s.logger.Error("error getting repo status by id", zap.Error(err))
		if errors.Is(err, pgx.ErrNoRows) {
			return response, errors.New(fmt.Sprintf("status id '%d' is not found", task.StatusID))
		}
		return response, constant.ErrInternalError
	}

	response.ID = task.ID
	response.Title = task.Title
	response.Description = task.Description
	response.Date = task.Date.Format(time.RFC3339)
	response.StatusName = status.Name
	response.CreatedAt = task.CreatedAt.Format(time.RFC3339)

	return response, nil
}
