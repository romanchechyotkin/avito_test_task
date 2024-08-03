package app

import (
	"github.com/gin-gonic/gin"
	"github.com/romanchechyotkin/avito_test_task/internal/controller/v1"
	"github.com/romanchechyotkin/avito_test_task/internal/service"
	"github.com/romanchechyotkin/avito_test_task/pkg/httpsrv"
	"github.com/romanchechyotkin/avito_test_task/pkg/postgresql"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/romanchechyotkin/avito_test_task/internal/config"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"
	"github.com/romanchechyotkin/avito_test_task/pkg/migrations"
	"github.com/romanchechyotkin/avito_test_task/schema"
)

func Run() {
	log := logger.New()
	log.Debug("app starting")

	cfg, err := config.New(log)
	if err != nil {
		log.Error("failed to init config", logger.Error(err))
	}

	log.Debug("app configuration", slog.Any("cfg", cfg))

	log.Debug("migrations starting")
	migrations.Migrate(log, &schema.DB, "postgres://postgres:5432@localhost:5432/estate_service?sslmode=disable")

	log.Debug("postgresql starting")
	postgres, err := postgresql.New(log, &cfg.Postgresql)
	if err != nil {
		log.Error("failed to init postgtresql", logger.Error(err))
		os.Exit(1)
	}
	postgres.Close()

	log.Debug("services init")
	services := service.NewServices()

	router := gin.Default()
	v1.NewRouter(log, router, services)

	log.Debug("server starting")
	server, err := httpsrv.New(log, cfg, router)

	// Waiting signal
	log.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-server.Notify():
		log.Error("app - Run - httpServer.Notify", logger.Error(err))
	}
}
