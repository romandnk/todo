package taskservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/romandnk/todo/internal/constant"
	"github.com/romandnk/todo/internal/entity"
	mock_storage "github.com/romandnk/todo/internal/repo/mock"
	mock_logger "github.com/romandnk/todo/pkg/logger/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"testing"
	"time"
)

const (
	errorLogger = "error"
	infoLogger  = "info"
)

func TestTaskService_CreateTask(t *testing.T) {
	date := time.Date(2024, 12, 07, 20, 49, 18, 0, time.UTC)

	type logger func(mock *mock_logger.MockLogger, msg string, args ...any)
	type createTask func(mock *mock_storage.MockTask, ctx context.Context, task entity.Task, expectedID int, expectedError error)
	type getStatusByName func(mock *mock_storage.MockStatus, ctx context.Context, name string, expectedStatus entity.Status, expectedError error)

	testCases := []struct {
		name                string
		input               CreateTaskParams
		loggerMock          logger
		loggerMsg           string
		loggerArgs          []any
		statusMock          getStatusByName
		taskMock            createTask
		expectedStatusName  string
		expectedStatus      entity.Status
		expectedStatusError error
		expectedTask        entity.Task
		expectedTaskID      int
		expectedTaskError   error
		expectedOutput      CreateTaskResponse
		expectedError       error
	}{
		{
			name: "OK",
			input: CreateTaskParams{
				Title:       "Test",
				Description: "Test",
				StatusName:  "Выполнено",
				Date:        "2024-12-07T20:49:18Z",
			},
			statusMock: func(mock *mock_storage.MockStatus, ctx context.Context, name string, expectedStatus entity.Status, expectedError error) {
				mock.EXPECT().GetStatusByName(ctx, name).Return(expectedStatus, expectedError)
			},
			taskMock: func(mock *mock_storage.MockTask, ctx context.Context, task entity.Task, expectedID int, expectedError error) {
				mock.EXPECT().CreateTask(ctx, task).Return(expectedID, expectedError)
			},
			expectedStatusName: "выполнено",
			expectedStatus: entity.Status{
				ID:   1,
				Name: "выполнено",
			},
			expectedTask: entity.Task{
				Title:       "Test",
				Description: "Test",
				StatusID:    1,
				Date:        date,
				Deleted:     false,
				DeletedAt:   time.Time{},
			},
			expectedTaskID: 1,
			expectedOutput: CreateTaskResponse{
				ID: 1,
			},
		},
		{
			name: "title is empty",
			input: CreateTaskParams{
				Title:       "",
				Description: "Test",
				StatusName:  "Выполнено",
				Date:        "2024-12-07T20:49:18Z",
			},
			expectedError: constant.ErrEmptyTitle,
		},
		{
			name: "OK",
			input: CreateTaskParams{
				Title:       "Test",
				Description: "Test",
				StatusName:  "test",
				Date:        "2024-12-07T20:49:18Z",
			},
			loggerMock: func(mock *mock_logger.MockLogger, msg string, args ...any) {
				mock.EXPECT().Error(msg, args...)
			},
			loggerMsg:  "error getting repo status by name",
			loggerArgs: []any{zap.Error(pgx.ErrNoRows)},
			statusMock: func(mock *mock_storage.MockStatus, ctx context.Context, name string, expectedStatus entity.Status, expectedError error) {
				mock.EXPECT().GetStatusByName(ctx, name).Return(expectedStatus, expectedError)
			},
			expectedStatusName:  "test",
			expectedStatusError: pgx.ErrNoRows,
			expectedError:       errors.New(fmt.Sprintf("status name '%s' is not found", "test")),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			taskStorage := mock_storage.NewMockTask(ctrl)
			statusStorage := mock_storage.NewMockStatus(ctrl)
			log := mock_logger.NewMockLogger(ctrl)

			taskService := NewTaskService(taskStorage, statusStorage, log)

			if tc.statusMock != nil {
				tc.statusMock(statusStorage, ctx, tc.expectedStatusName, tc.expectedStatus, tc.expectedStatusError)
			}

			if tc.taskMock != nil {
				tc.taskMock(taskStorage, ctx, tc.expectedTask, tc.expectedTaskID, tc.expectedTaskError)
			}

			if tc.loggerMock != nil {
				tc.loggerMock(log, tc.loggerMsg, tc.loggerArgs...)
			}

			output, err := taskService.CreateTask(ctx, tc.input)
			if err != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedOutput, output)
		})
	}
}
