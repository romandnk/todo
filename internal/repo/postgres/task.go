package postgresrepo

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/romandnk/todo/internal/constant"
	"github.com/romandnk/todo/internal/entity"
	postgres "github.com/romandnk/todo/pkg/storage"
	"github.com/romandnk/todo/pkg/utils"
	"time"
)

type TaskRepo struct {
	db postgres.PgxPool
}

func NewTaskRepo(db postgres.PgxPool) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) CreateTask(ctx context.Context, task entity.Task) (int, error) {
	var id int

	values := []any{task.Title, task.Description, task.StatusID, task.Date, task.Deleted, task.CreatedAt, task.DeletedAt}
	query := fmt.Sprintf(`
		INSERT INTO %[1]s
		(title, description, status_id, date, deleted, created_at, deleted_at)
		VALUES %[2]s
		RETURNING id
	`, constant.TasksTable, utils.SetPlaceholders(constant.PlaceholderDollar, len(values)))

	err := pgxscan.Get(ctx, r.db, &id, query, values...)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (r *TaskRepo) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	var task entity.Task

	query := fmt.Sprintf(`
		SELECT 
		    t.id, 
		    t.title, 
		    t.description, 
		    s.name, 
		    t.date, 
		    t.deleted, 
		    t.created_at, 
		    t.deleted_at
		FROM %[1]s t
		JOIN %[2]s s ON t.status_id=s.id
		WHERE t.id=$1 AND deleted=false
	`, constant.TasksTable, constant.StatusesTable)

	err := pgxscan.Get(ctx, r.db, &task, query, id)
	if err != nil {
		return task, err
	}

	return task, nil
}

func (r *TaskRepo) GetTasksByTitle(ctx context.Context, title string) ([]*entity.Task, error) {
	var tasks []*entity.Task

	query := fmt.Sprintf(`
		SELECT 
		    t.id, 
		    t.title, 
		    t.description, 
		    s.name, 
		    t.date, 
		    t.deleted, 
		    t.created_at, 
		    t.deleted_at
		FROM %[1]s t
		JOIN %[2]s s ON t.status_id=s.id
		WHERE t.title LIKE '%$1%' AND deleted=false
	`, constant.TasksTable, constant.StatusesTable)

	err := pgxscan.Select(ctx, r.db, &tasks, query, title)
	if err != nil {
		return tasks, err
	}

	return tasks, nil
}

func (r *TaskRepo) UpdateTaskByID(ctx context.Context, id int, task entity.Task) error {
	newTask := utils.CheckEmptyTaskFields(task)

	values := []any{newTask.Title, newTask.Description, newTask.StatusID, newTask.Date, id}
	query := fmt.Sprintf(`
		UPDATE %[1]s
		SET 
			title=COALESCE($1, title),
			description=COALESCE($2, description),
			status_id=COALESCE($3, status_id),
			date=COALESCE($4, date)
		WHERE id=$5 AND deleted=false
	`, constant.TasksTable)

	_, err := r.db.Exec(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepo) DeleteTaskByID(ctx context.Context, id int) error {
	now := time.Now().UTC()

	query := fmt.Sprintf(`
		UPDATE %[1]s
		SET 
		    deleted=true,
		    deleted_at=$1
		WHERE id=$2 AND deleted=false
	`, constant.TasksTable)

	res, err := r.db.Exec(ctx, query, now, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return constant.ErrTaskIDNotExists
	}

	return nil
}
