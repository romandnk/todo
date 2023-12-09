package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	mock_service "github.com/romandnk/todo/internal/service/mock"
	taskservice "github.com/romandnk/todo/internal/service/task"
	mock_logger "github.com/romandnk/todo/pkg/logger/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTaskRoutes_CreateTask(t *testing.T) {
	url := "/api/v1/tasks"

	type argsTask struct {
		input         taskservice.CreateTaskParams
		output        taskservice.CreateTaskResponse
		expectedError error
	}

	type argsLogger struct {
		msg    string
		fields []any
	}

	type mockTaskBehaviour func(m *mock_service.MockTask, args argsTask)
	type mockLoggerBehaviour func(m *mock_logger.MockLogger, args argsLogger)

	testCases := []struct {
		name                 string
		argsTask             argsTask
		argsLogger           argsLogger
		taskM                mockTaskBehaviour
		loggerM              mockLoggerBehaviour
		requestBody          map[string]interface{}
		expectedResponseBody string
		expectedHTTPCode     int
	}{
		{
			name: "OK",
			argsTask: argsTask{
				input: taskservice.CreateTaskParams{
					Title:       "Test",
					Description: "Test",
					StatusName:  "done",
					Date:        "2024-12-07T20:49:18Z",
				},
				output: taskservice.CreateTaskResponse{ID: 1},
			},
			taskM: func(m *mock_service.MockTask, args argsTask) {
				m.EXPECT().CreateTask(gomock.Any(), args.input).Return(args.output, args.expectedError)
			},
			requestBody: map[string]interface{}{
				"title":       "Test",
				"description": "Test",
				"status_name": "done",
				"date":        "2024-12-07T20:49:18Z",
			},
			expectedResponseBody: `{"id":1}`,
			expectedHTTPCode:     http.StatusCreated,
		},
		{
			name: "error validating json body fields (empty title)",
			argsLogger: argsLogger{
				msg:    "error binding json body",
				fields: []any{zap.String("error", "Key: 'CreateTaskParams.Title' Error:Field validation for 'Title' failed on the 'required' tag")},
			},
			loggerM: func(m *mock_logger.MockLogger, args argsLogger) {
				m.EXPECT().Error(args.msg, args.fields)
			},
			requestBody: map[string]interface{}{
				"description": "Test",
				"status_name": "done",
				"date":        "2024-12-07T20:49:18Z",
			},
			expectedResponseBody: `{"message":"error binding json body","error":"Key: 'CreateTaskParams.Title' Error:Field validation for 'Title' failed on the 'required' tag"}`,
			expectedHTTPCode:     http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			taskService := mock_service.NewMockTask(ctrl)
			logger := mock_logger.NewMockLogger(ctrl)

			if tc.taskM != nil {
				tc.taskM(taskService, tc.argsTask)
			}

			if tc.loggerM != nil {
				tc.loggerM(logger, tc.argsLogger)
			}

			taskR := taskRoutes{
				task:   taskService,
				logger: logger,
			}

			r := gin.Default()
			r.POST(url, taskR.CreateTask)

			jsonBody, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			require.Equal(t, tc.expectedHTTPCode, w.Code)

			require.Equal(t, []byte(tc.expectedResponseBody), w.Body.Bytes())
		})
	}
}
