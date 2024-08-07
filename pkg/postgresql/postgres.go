package postgresql

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/romanchechyotkin/avito_test_task/internal/config"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxPool interface {
	Close()
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Ping(ctx context.Context) error
}

type Postgres struct {
	log  *slog.Logger
	Pool PgxPool
}

func New(log *slog.Logger, cfg *config.Postgresql) (*Postgres, error) {
	log = log.With(slog.String("component", "postgresql"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
	)

	log.Debug("connection url", slog.String("url", url))

	if cfg.AutoCreate {
		exists, err := CreateDatabase(ctx, url)
		if err != nil {
			log.Error("failed to create database", logger.Error(err))

			return nil, err
		}
		if exists {
			log.Info("the database already exists")
		} else {
			log.Info("the database was created successfully")
		}
	}

	pgCfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Error("failed to parse postgres config", logger.Error(err))

		return nil, err
	}

	conn, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		log.Error("pool constructor error", logger.Error(err))

		return nil, err
	}

	return &Postgres{
		log:  log,
		Pool: conn,
	}, conn.Ping(ctx)
}

func (p *Postgres) Close() {
	p.log.Debug("close postgres connection")
	if p.Pool != nil {
		p.Pool.Close()
	}
}
