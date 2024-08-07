package utest

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/romanchechyotkin/avito_test_task/internal/config"
	"github.com/romanchechyotkin/avito_test_task/pkg/migrations"
	"github.com/romanchechyotkin/avito_test_task/schema"

	"github.com/romanchechyotkin/avito_test_task/pkg/logger"
	"github.com/romanchechyotkin/avito_test_task/pkg/postgresql"
)

func Prepare() (*slog.Logger, *config.Config, *postgresql.Postgres, error) {
	log := logger.New()
	cfg, err := config.New(log)
	if err != nil {
		return nil, nil, nil, err
	}

	pg, err := postgresql.New(log, &cfg.Postgresql)
	if err != nil {
		return nil, nil, nil, err
	}

	migrations.Migrate(log, &schema.DB, &cfg.Postgresql)

	return log, cfg, pg, nil
}

func TeardownTable(log *slog.Logger, pg *postgresql.Postgres, tableName string) {
	exec, err := pg.Pool.Exec(context.Background(), fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName))
	if err != nil {
		log.Error("failed to truncate table", slog.String("table", tableName), slog.String("query", "TRUNCATE"), logger.Error(err))
		return
	}
	log.Debug("truncated table", slog.String("table", tableName), slog.Int64("rows affected", exec.RowsAffected()))
}
