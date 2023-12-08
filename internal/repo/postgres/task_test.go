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

func TestTaskRepoCreateTask(t *testing.T) {
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

func TestTaskRepo_GetAllTasks(t *testing.T) {
	now := time.Now().UTC()

	testCases := []struct {
		name          string
		query         string
		statusID      int
		limit         int
		lastID        int
		date          time.Time
		args          []any
		expectedTasks []*entity.Task
		expectedError error
	}{
		{
			name: "OK with all input data",
			query: `
				SELECT 
		    		id, 
		    		title, 
		    		description, 
		    		status_id, 
		    		date,  
		    		created_at
				FROM tasks
				WHERE deleted=false AND status_id=$1 AND date BETWEEN $2 AND $3 AND id>$4
				ORDER BY id
				LIMIT $5
			`,
			statusID: 1,
			limit:    5,
			lastID:   1,
			date:     now,
			args: []any{
				1,
				time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
				time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).
					AddDate(0, 0, 1).Add(-time.Nanosecond),
				1,
				5,
			},
			expectedTasks: []*entity.Task{
				{
					ID:          1,
					Title:       "Test",
					Description: "Test",
					StatusID:    1,
					Date:        now,
					Deleted:     false,
					CreatedAt:   now,
					DeletedAt:   time.Time{},
				},
				{
					ID:          2,
					Title:       "Test",
					Description: "Test",
					StatusID:    1,
					Date:        now,
					Deleted:     false,
					CreatedAt:   now,
					DeletedAt:   time.Time{},
				},
			},
			expectedError: nil,
		},
		{
			name: "OK with status id and date",
			query: `
				SELECT 
		    		id, 
		    		title, 
		    		description, 
		    		status_id, 
		    		date,  
		    		created_at
				FROM tasks
				WHERE deleted=false AND status_id=$1 AND date BETWEEN $2 AND $3 AND id>$4
				ORDER BY id
			`,
			statusID: 1,
			limit:    0,
			lastID:   1,
			date:     now,
			args: []any{
				1,
				time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
				time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).
					AddDate(0, 0, 1).Add(-time.Nanosecond),
				1,
			},
			expectedTasks: []*entity.Task{
				{
					ID:          1,
					Title:       "Test",
					Description: "Test",
					StatusID:    1,
					Date:        now,
					Deleted:     false,
					CreatedAt:   now,
					DeletedAt:   time.Time{},
				},
				{
					ID:          2,
					Title:       "Test",
					Description: "Test",
					StatusID:    1,
					Date:        now,
					Deleted:     false,
					CreatedAt:   now,
					DeletedAt:   time.Time{},
				},
			},
			expectedError: nil,
		},
		{
			name: "OK without input data",
			query: `
				SELECT 
		    		id, 
		    		title, 
		    		description, 
		    		status_id, 
		    		date,  
		    		created_at
				FROM tasks
				WHERE deleted=false AND id>$1
				ORDER BY id
			`,
			statusID: 0,
			limit:    0,
			lastID:   0,
			date:     time.Time{},
			args:     []any{0},
			expectedTasks: []*entity.Task{
				{
					ID:          1,
					Title:       "Test",
					Description: "Test",
					StatusID:    1,
					Date:        now,
					Deleted:     false,
					CreatedAt:   now,
					DeletedAt:   time.Time{},
				},
				{
					ID:          2,
					Title:       "Test",
					Description: "Test",
					StatusID:    1,
					Date:        now,
					Deleted:     false,
					CreatedAt:   now,
					DeletedAt:   time.Time{},
				},
			},
			expectedError: nil,
		},
		{
			name: "No rows in result set",
			query: `
				SELECT 
		    		id, 
		    		title, 
		    		description, 
		    		status_id, 
		    		date,  
		    		created_at
				FROM tasks
				WHERE deleted=false AND id>$1
				ORDER BY id
			`,
			statusID:      0,
			limit:         0,
			lastID:        0,
			date:          time.Time{},
			args:          []any{0},
			expectedTasks: []*entity.Task{},
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

			columns := []string{"id", "title", "description", "status_id", "date", "created_at"}
			rows := pgxmock.NewRows(columns)
			for _, task := range tc.expectedTasks {
				rows.AddRow(
					task.ID,
					task.Title,
					task.Description,
					task.StatusID,
					task.Date,
					task.CreatedAt,
				)
			}

			if tc.expectedError == nil {
				mock.ExpectQuery(regexp.QuoteMeta(tc.query)).WithArgs(tc.args...).WillReturnRows(rows)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(tc.query)).WithArgs(tc.args...).WillReturnError(tc.expectedError)
			}

			storage := NewTaskRepo(mock)

			tasks, err := storage.GetAllTasks(ctx, tc.statusID, tc.limit, tc.lastID, tc.date)
			require.ErrorIs(t, err, tc.expectedError)
			require.ElementsMatch(t, tc.expectedTasks, tasks)

			require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
		})
	}
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
		    		created_at
				FROM %[1]s
				WHERE id=$1
			`, constant.TasksTable)

			columns := []string{"id", "title", "description", "status_id", "date", "created_at"}
			rows := pgxmock.NewRows(columns).
				AddRow(
					tc.expectedTask.ID,
					tc.expectedTask.Title,
					tc.expectedTask.Description,
					tc.expectedTask.StatusID,
					tc.expectedTask.Date,
					tc.expectedTask.CreatedAt,
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
