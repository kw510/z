package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pg *pgxpool.Pool

func Init(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("failed to create pool: %w", err)
	}
	err = pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping pool: %w", err)
	}

	pg = pool

	return nil
}

func Tx(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error {
	tx, err := pg.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := fn(ctx, tx); err != nil {
		return fmt.Errorf("failed to execute transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func Close() {
	pg.Close()
}
