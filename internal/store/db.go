package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewPostgres(databaseURL string) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}
	return &DB{pool: pool}, nil
}

func (d *DB) Close() {
	if d.pool != nil {
		d.pool.Close()
	}
}

// Simple ping helper
func (d *DB) Ping(ctx context.Context) error {
	if d.pool == nil {
		return errors.New("db not initialized")
	}
	return d.pool.Ping(ctx)
}
