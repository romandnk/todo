package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/romandnk/todo/config"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxPool interface {
	Close()
	Acquire(ctx context.Context) (*pgxpool.Pool, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Ping(ctx context.Context) error
}

func NewStorage(ctx context.Context, cfg config.Postgres) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	pgxConf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pgxConf.MaxConns = cfg.MaxConns
	pgxConf.MinConns = cfg.MinConns

	db, err := pgxpool.NewWithConfig(ctx, pgxConf)
	if err != nil {
		return nil, fmt.Errorf("error creating new pgx pool: %w", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("error connecting pgx pool: %w", err)
	}

	return db, nil
}
