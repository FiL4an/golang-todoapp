package core_pgx_pool

import (
	"context"
	"fmt"
	"time"

	core_postgres_pool "github.com/FiL4an/golang-todoapp/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

func NewConnectionPool(config Config, ctx context.Context) (*Pool, error) {
	connections := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database)

	pgxconfig, err := pgxpool.ParseConfig(connections)
	if err != nil {
		return &Pool{}, fmt.Errorf("parse pgxconfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pgxpool ping: %w", err)
	}

	return &Pool{
		Pool:      pool,
		opTimeout: config.Timeout,
	}, nil
}

func (p *Pool) OpTimeout() time.Duration {
	return p.opTimeout
}

func (p *Pool) Query(
	ctx context.Context,
	sql string,
	args ...any) (core_postgres_pool.Rows, error) {

	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return pgxRows{rows}, nil
}
func (p *Pool) QueryRow(
	ctx context.Context,
	sql string,
	args ...any) core_postgres_pool.Row {
	row := p.Pool.QueryRow(ctx, sql, args...)

	return pgxRow{row}
}
func (p *Pool) Exec(
	ctx context.Context,
	sql string,
	arguments ...any) (core_postgres_pool.CommandTag, error) {

	commandTag, err := p.Pool.Exec(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}

	return pgxCommandTagSrcuct{CommandTag: commandTag}, nil

}

func (p *Pool) OpTimeoit() time.Duration {
	return p.opTimeout
}
