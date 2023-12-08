package postgresrepo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/romandnk/todo/internal/constant"
	"github.com/romandnk/todo/internal/entity"
	"github.com/romandnk/todo/pkg/utils"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

func TestTaskRepo_CreateTask(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	ctx := context.Background()

	now := time.Now().UTC()
	inputTask := entity.Task{
		Title:       "Test",
		Description: "Test",
		StatusID:    1,
		Date:        now,
		Deleted:     false,
		CreatedAt:   now,
		DeletedAt:   time.Time{},
	}
	expectedID := 1

	query := fmt.Sprintf(`
		INSERT INTO %[1]s
		(title, description, status_id, date, deleted, created_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, constant.TasksTable)

	columns := []string{"id"}
	rows := pgxmock.NewRows(columns).AddRow(expectedID)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(
		inputTask.Title,
		inputTask.Description,
		inputTask.StatusID,
		inputTask.Date,
		inputTask.Deleted,
		inputTask.CreatedAt,
		inputTask.DeletedAt,
	).WillReturnRows(rows)

	storage := NewTaskRepo(mock)

	id, err := storage.CreateTask(ctx, inputTask)
	require.NoError(t, err)
	require.Equal(t, expectedID, id)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestTaskRepo_GetTaskByID(t *testing.T) {
	now := time.Now().UTC()

	testCases := []struct {
		name          string
		expectedID    int
		expectedTask  entity.Task
		expectedError error
	}{
		{
			name:       "OK",
			expectedID: 1,
			expectedTask: entity.Task{
				ID:          1,
				Title:       "Test",
				Description: "Test",
				StatusID:    1,
				Date:        now,
				Deleted:     false,
				CreatedAt:   now,
				DeletedAt:   time.Time{},
			},
			expectedError: nil,
		},
		{
			name:          "No rows in result set",
			expectedID:    1,
			expectedTask:  entity.Task{},
			expectedError: pgx.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			ctx := context.Background()

			query := fmt.Sprintf(`
				SELECT 
		    		id, 
		    		title, 
		    		description, 
		    		status_id, 
		    		date, 
		    		deleted, 
		    		created_at, 
		    		deleted_at
				FROM %[1]s
				WHERE id=$1
			`, constant.TasksTable)

			columns := []string{"id", "title", "description", "status_id", "date", "deleted", "created_at", "deleted_at"}
			rows := pgxmock.NewRows(columns).
				AddRow(
					tc.expectedTask.ID,
					tc.expectedTask.Title,
					tc.expectedTask.Description,
					tc.expectedTask.StatusID,
					tc.expectedTask.Date,
					tc.expectedTask.Deleted,
					tc.expectedTask.CreatedAt,
					tc.expectedTask.DeletedAt,
				)

			if tc.expectedError == nil {
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(tc.expectedID).WillReturnRows(rows)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(tc.expectedID).WillReturnError(tc.expectedError)
			}

			storage := NewTaskRepo(mock)

			task, err := storage.GetTaskByID(ctx, tc.expectedID)
			require.ErrorIs(t, err, tc.expectedError)
			require.Equal(t, tc.expectedTask, task)

			require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
		})
	}
}

func TestTaskRepo_DeleteTaskByID(t *testing.T) {
	update := "update"

	testCases := []struct {
		name          string
		expectedID    int
		expectedError error
	}{
		{
			name:          "OK",
			expectedID:    1,
			expectedError: nil,
		},
		{
			name:          "Task with id isn't found",
			expectedID:    1,
			expectedError: constant.ErrTaskIDNotExists,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			ctx := context.Background()

			query := fmt.Sprintf(`
				UPDATE %[1]s
				SET 
				    deleted=true,
				    deleted_at=$1
				WHERE id=$2 AND deleted=false
			`, constant.TasksTable)

			if tc.expectedError == nil {
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(pgxmock.AnyArg(), tc.expectedID).WillReturnResult(pgxmock.NewResult(update, 1))
			} else {
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(pgxmock.AnyArg(), tc.expectedID).WillReturnError(tc.expectedError)
			}

			storage := NewTaskRepo(mock)

			err = storage.DeleteTaskByID(ctx, tc.expectedID)
			require.ErrorIs(t, err, tc.expectedError)

			require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
		})
	}
}

func TestTaskRepo_UpdateTaskByID(t *testing.T) {
	now := time.Now().UTC()
	update := "update"

	testCases := []struct {
		name                string
		expectedID          int
		expectedUpdatedTask entity.Task
		expectedInput       utils.TaskToUpdate
		expectedError       error
	}{
		{
			name:       "OK full task for updating",
			expectedID: 1,
			expectedUpdatedTask: entity.Task{
				Title:       "test",
				Description: "test",
				StatusID:    2,
				Date:        now,
			},
			expectedInput: utils.TaskToUpdate{
				Title: sql.NullString{
					String: "test",
					Valid:  true,
				},
				Description: sql.NullString{
					String: "test",
					Valid:  true,
				},
				StatusID: sql.NullInt16{
					Int16: 2,
					Valid: true,
				},
				Date: sql.NullTime{
					Time:  now,
					Valid: true,
				},
			},
			expectedError: nil,
		},
		{
			name:       "OK with only updating title ans status id",
			expectedID: 1,
			expectedUpdatedTask: entity.Task{
				Title:    "test",
				StatusID: 2,
			},
			expectedInput: utils.TaskToUpdate{
				Title: sql.NullString{
					String: "test",
					Valid:  true,
				},
				Description: sql.NullString{
					String: "",
					Valid:  false,
				},
				StatusID: sql.NullInt16{
					Int16: 2,
					Valid: true,
				},
				Date: sql.NullTime{
					Time:  time.Time{},
					Valid: false,
				},
			},
			expectedError: nil,
		},
		{
			name:                "No rows in result set",
			expectedID:          1,
			expectedUpdatedTask: entity.Task{},
			expectedInput:       utils.TaskToUpdate{},
			expectedError:       constant.ErrTaskIDNotExists,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			ctx := context.Background()

			query := fmt.Sprintf(`
				UPDATE %[1]s
				SET 
					title=COALESCE($1, title),
					description=COALESCE($2, description),
					status_id=COALESCE($3, status_id),
					date=COALESCE($4, date)
				WHERE id=$5 AND deleted=false
			`, constant.TasksTable)

			if tc.expectedError == nil {
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
					tc.expectedInput.Title,
					tc.expectedInput.Description,
					tc.expectedInput.StatusID,
					pgxmock.AnyArg(),
					tc.expectedID,
				).WillReturnResult(pgxmock.NewResult(update, 1))
			} else {
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
					tc.expectedInput.Title,
					tc.expectedInput.Description,
					tc.expectedInput.StatusID,
					pgxmock.AnyArg(),
					tc.expectedID,
				).WillReturnError(tc.expectedError)
			}

			storage := NewTaskRepo(mock)

			err = storage.UpdateTaskByID(ctx, tc.expectedID, tc.expectedUpdatedTask)
			require.ErrorIs(t, err, tc.expectedError)

			require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
		})
	}
}
