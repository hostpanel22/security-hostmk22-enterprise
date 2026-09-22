package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	DefaultMaxOpenConns    = 25
	DefaultMaxIdleConns    = 10
	DefaultConnMaxLifetime = 30 * time.Minute
	DefaultPingTimeout     = 5 * time.Second
)

type Config struct {
	DatabaseURL  string
	MaxOpenConns int
	MaxIdleConns int
	ConnMaxLife  time.Duration
	PingTimeout  time.Duration
}

type Postgres struct {
	db *sql.DB
}

func New(cfg Config) (*Postgres, error) {
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("database url is required")
	}

	if cfg.MaxOpenConns <= 0 {
		cfg.MaxOpenConns = DefaultMaxOpenConns
	}

	if cfg.MaxIdleConns <= 0 {
		cfg.MaxIdleConns = DefaultMaxIdleConns
	}

	if cfg.ConnMaxLife <= 0 {
		cfg.ConnMaxLife = DefaultConnMaxLifetime
	}

	if cfg.PingTimeout <= 0 {
		cfg.PingTimeout = DefaultPingTimeout
	}

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLife)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		cfg.PingTimeout,
	)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &Postgres{
		db: db,
	}, nil
}

func (p *Postgres) DB() *sql.DB {
	return p.db
}

func (p *Postgres) Close() error {
	if p == nil || p.db == nil {
		return nil
	}

	return p.db.Close()
}

func (p *Postgres) HealthCheck() error {
	if p == nil || p.db == nil {
		return fmt.Errorf("database not initialized")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		DefaultPingTimeout,
	)
	defer cancel()

	return p.db.PingContext(ctx)
}

func (p *Postgres) Stats() sql.DBStats {
	if p == nil || p.db == nil {
		return sql.DBStats{}
	}

	return p.db.Stats()
}

func (p *Postgres) BeginTx(
	ctx context.Context,
	opts *sql.TxOptions,
) (*sql.Tx, error) {
	if p == nil || p.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	return p.db.BeginTx(ctx, opts)
}

func (p *Postgres) Exec(
	ctx context.Context,
	query string,
	args ...any,
) (sql.Result, error) {
	if p == nil || p.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	return p.db.ExecContext(ctx, query, args...)
}

func (p *Postgres) Query(
	ctx context.Context,
	query string,
	args ...any,
) (*sql.Rows, error) {
	if p == nil || p.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	return p.db.QueryContext(ctx, query, args...)
}

func (p *Postgres) QueryRow(
	ctx context.Context,
	query string,
	args ...any,
) *sql.Row {
	return p.db.QueryRowContext(ctx, query, args...)
}
