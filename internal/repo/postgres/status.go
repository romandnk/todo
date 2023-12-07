package postgresrepo

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/romandnk/todo/internal/constant"
	"github.com/romandnk/todo/internal/entity"
	postgres "github.com/romandnk/todo/pkg/storage"
	"github.com/romandnk/todo/pkg/utils"
)

type StatusRepo struct {
	db postgres.PgxPool
}

func NewStatusRepo(db postgres.PgxPool) *StatusRepo {
	return &StatusRepo{db: db}
}

func (r *StatusRepo) CreateStatus(ctx context.Context, status entity.Status) (int, error) {
	var id int

	values := []any{status.Name}
	query := fmt.Sprintf(`
		INSERT INTO %[1]s
		(name)
		VALUES %[2]s
		RETURNING id
	`, constant.StatusesTable, utils.SetPlaceholders(constant.PlaceholderDollar, len(values)))

	err := pgxscan.Get(ctx, r.db, &id, query, values...)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (r *StatusRepo) GetStatusByName(ctx context.Context, name string) (entity.Status, error) {
	var status entity.Status

	query := fmt.Sprintf(`
		SELECT id, name
		FROM %[1]s
		WHERE name=$1
	`, constant.StatusesTable)

	err := pgxscan.Get(ctx, r.db, &status, query, name)
	if err != nil {
		return status, err
	}

	return status, nil
}

func (r *StatusRepo) GetStatusByID(ctx context.Context, id int) (entity.Status, error) {
	var status entity.Status

	query := fmt.Sprintf(`
		SELECT id, name
		FROM %[1]s
		WHERE id=$1
	`, constant.StatusesTable)

	err := pgxscan.Get(ctx, r.db, &status, query, id)
	if err != nil {
		return status, err
	}

	return status, nil
}
