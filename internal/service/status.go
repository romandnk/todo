package service

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/romandnk/todo/internal/constant"
	"github.com/romandnk/todo/internal/entity"
	storage "github.com/romandnk/todo/internal/repo"
	"github.com/romandnk/todo/pkg/logger"
	"go.uber.org/zap"
	"strings"
	"unicode/utf8"
)

type statusService struct {
	status storage.Status
	logger logger.Logger
}

func newStatusService(status storage.Status, logger logger.Logger) *statusService {
	return &statusService{
		status: status,
		logger: logger,
	}
}

func (s *statusService) CreateStatus(ctx context.Context, params CreateStatusParams) (CreateStatusResponse, error) {
	var response CreateStatusResponse

	params.Name = strings.TrimSpace(params.Name)

	if params.Name == "" {
		return response, constant.ErrEmptyStatusName
	}

	if utf8.RuneCountInString(params.Name) > 36 {
		return response, constant.ErrTooLongStatusName
	}

	status := entity.Status{
		Name: params.Name,
	}
	id, err := s.status.CreateStatus(ctx, status)
	if err != nil {
		s.logger.Error("error creating repo task", zap.Error(err))
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				if pgErr.ColumnName == "name" {
					return response, constant.ErrStatusNameExists
				}
			}
		}

		return response, constant.ErrInternalError
	}

	response.ID = id

	return response, nil
}
