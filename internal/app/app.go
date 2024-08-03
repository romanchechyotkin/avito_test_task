package app

import (
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"
	"github.com/romanchechyotkin/avito_test_task/pkg/migrations"
	"github.com/romanchechyotkin/avito_test_task/schema"
)

func Run() {
	log := logger.New()
	log.Debug("app starting")

	log.Info("migrations starting")
	migrations.Migrate(log, &schema.DB, "postgres://postgres:5432@localhost:5432/estate_service?sslmode=disable")

}
