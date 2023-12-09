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
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTaskRoutes_CreateTask(t *testing.T) {
	url := "/api/v1/tasks"

	type args struct {
		input         taskservice.CreateTaskParams
		output        taskservice.CreateTaskResponse
		expectedError error
	}

	type mockBehaviour func(m *mock_service.MockTask, args args)

	testCases := []struct {
		name                 string
		args                 args
		m                    mockBehaviour
		requestBody          map[string]interface{}
		expectedResponseBody map[string]interface{}
	}{
		{
			name: "OK",
			args: args{
				input: taskservice.CreateTaskParams{
					Title:       "Test",
					Description: "Test",
					StatusName:  "done",
					Date:        "2024-12-07T20:49:18Z",
				},
				output: taskservice.CreateTaskResponse{ID: 1},
			},
			m: func(m *mock_service.MockTask, args args) {
				m.EXPECT().CreateTask(gomock.Any(), args.input).Return(args.output, args.expectedError)
			},
			requestBody: map[string]interface{}{
				"title":       "Test",
				"description": "Test",
				"status_name": "done",
				"date":        "2024-12-07T20:49:18Z",
			},
			expectedResponseBody: map[string]interface{}{
				"id": 1,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			taskService := mock_service.NewMockTask(ctrl)
			logger := mock_logger.NewMockLogger(ctrl)

			tc.m(taskService, tc.args)

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

			require.Equal(t, http.StatusCreated, w.Code)

			expectedResponseBody, err := json.Marshal(tc.expectedResponseBody)
			require.NoError(t, err)

			require.Equal(t, []byte(expectedResponseBody), w.Body.Bytes())
		})
	}
}
