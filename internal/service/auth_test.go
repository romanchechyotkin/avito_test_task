//go:build integration

package service

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/romanchechyotkin/avito_test_task/internal/config"
	"github.com/romanchechyotkin/avito_test_task/internal/repo"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"
	"github.com/romanchechyotkin/avito_test_task/pkg/migrations"
	"github.com/romanchechyotkin/avito_test_task/pkg/postgresql"
	"github.com/romanchechyotkin/avito_test_task/schema"

	"github.com/stretchr/testify/require"
)

func TestAuthService_CreateUser(t *testing.T) {
	log := logger.New()

	cfg, err := config.New(log)
	require.NoError(t, err)

	log.Debug("app configuration", slog.Any("cfg", cfg.Postgresql))

	pg, err := postgresql.New(log, &cfg.Postgresql)
	require.NoError(t, err)

	defer func() {
		exec, err := pg.Pool.Exec(context.Background(), "TRUNCATE TABLE users CASCADE")
		if err != nil {
			log.Error("failed to truncate users table", logger.Error(err))
			return
		}
		log.Debug("truncated users table", slog.Any("exec", exec))
	}()

	err = migrations.Migrate(log, &schema.DB, fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Postgresql.User,
		cfg.Postgresql.Password,
		cfg.Postgresql.Host,
		cfg.Postgresql.Port,
		cfg.Postgresql.Database,
		cfg.Postgresql.SSLMode,
	))
	require.NoError(t, err)

	repositories := repo.NewRepositories(log, pg)

	authService := NewAuthService(log, repositories.User, cfg.JWT.SignKey, cfg.JWT.TokenTTL)

	t.Run("create user", func(t *testing.T) {
		log.Debug("creating invalid user")
		userID, err := authService.CreateUser(context.Background(), &AuthCreateUserInput{
			Email:    "test",
			Password: "test",
			UserType: "test",
		})
		require.Error(t, err)
		require.Equal(t, "", userID)

		log.Debug("creating correct user")
		userID, err = authService.CreateUser(context.Background(), &AuthCreateUserInput{
			Email:    "test",
			Password: "test",
			UserType: "moderator",
		})
		require.NoError(t, err)
		require.True(t, len(userID) > 0)

		log.Debug("creating existing user")
		userID, err = authService.CreateUser(context.Background(), &AuthCreateUserInput{
			Email:    "test",
			Password: "test",
			UserType: "moderator",
		})
		require.ErrorIs(t, err, ErrUserExists)
		require.Equal(t, "", userID)
	})

	t.Run("create user with empty password", func(t *testing.T) {
		log.Debug("creating user with empty password")
		userID, err := authService.CreateUser(context.Background(), &AuthCreateUserInput{
			Email:    "test",
			Password: "",
			UserType: "moderator",
		})
		require.Error(t, err)
		require.Equal(t, "", userID)
	})
}

func TestAuthService_GenerateToken(t *testing.T) {

}
