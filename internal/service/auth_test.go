//go:build integration

package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/romanchechyotkin/avito_test_task/internal/config"
	"github.com/romanchechyotkin/avito_test_task/internal/repo"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"
	"github.com/romanchechyotkin/avito_test_task/pkg/migrations"
	"github.com/romanchechyotkin/avito_test_task/pkg/postgresql"
	"github.com/romanchechyotkin/avito_test_task/schema"

	"github.com/stretchr/testify/assert"
)

func TestAuthService_CreateUser(t *testing.T) {
	log := logger.NewDiscardLogger()

	cfg, err := config.New(log)
	assert.NoError(t, err)

	pg, err := postgresql.New(log, &cfg.Postgresql)
	assert.NoError(t, err)
	defer pg.Close()

	migrations.Migrate(log, &schema.DB, fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Postgresql.User,
		cfg.Postgresql.Password,
		cfg.Postgresql.Host,
		cfg.Postgresql.Port,
		cfg.Postgresql.Database,
		cfg.Postgresql.SSLMode,
	))

	repositories := repo.NewRepositories(log, pg)

	authService := NewAuthService(log, repositories.User, cfg.JWT.SignKey, cfg.JWT.TokenTTL)

	_, err = authService.CreateUser(context.Background(), &AuthCreateUserInput{
		Email:    "test",
		Password: "test",
		UserType: "test",
	})
	assert.Error(t, err)

	_, err = authService.CreateUser(context.Background(), &AuthCreateUserInput{
		Email:    "test",
		Password: "test",
		UserType: "moderator",
	})
	assert.NoError(t, err)

	_, err = authService.CreateUser(context.Background(), &AuthCreateUserInput{
		Email:    "test",
		Password: "test",
		UserType: "moderator",
	})
	assert.ErrorIs(t, err, ErrUserExists)
}
