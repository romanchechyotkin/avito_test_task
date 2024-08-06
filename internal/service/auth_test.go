//go:build integration

package service

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/romanchechyotkin/avito_test_task/internal/config"
	"github.com/romanchechyotkin/avito_test_task/internal/entity"
	"github.com/romanchechyotkin/avito_test_task/internal/repo"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"
	"github.com/romanchechyotkin/avito_test_task/pkg/migrations"
	"github.com/romanchechyotkin/avito_test_task/pkg/postgresql"
	"github.com/romanchechyotkin/avito_test_task/schema"

	"github.com/stretchr/testify/require"
)

var (
	log    = logger.New()
	cfg, _ = config.New(log)
)

func TestAuthService_CreateUser(t *testing.T) {
	log.Debug("test configuration", slog.Any("cfg", cfg.Postgresql))

	pg, err := postgresql.New(log, &cfg.Postgresql)
	require.NoError(t, err)

	defer func() {
		exec, err := pg.Pool.Exec(context.Background(), "TRUNCATE TABLE users CASCADE")
		if err != nil {
			log.Error("failed to truncate users table", logger.Error(err))
			return
		}
		log.Debug("truncated users table", slog.Int64("rows affected", exec.RowsAffected()))
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

func TestAuthService_GenerateToken_ParseToken(t *testing.T) {
	log.Debug("test configuration", slog.Any("cfg", cfg.Postgresql))

	pg, err := postgresql.New(log, &cfg.Postgresql)
	require.NoError(t, err)

	defer func() {
		exec, err := pg.Pool.Exec(context.Background(), "TRUNCATE TABLE users CASCADE")
		if err != nil {
			log.Error("failed to truncate users table", logger.Error(err))
			return
		}
		log.Debug("truncated users table", slog.Int64("rows affected", exec.RowsAffected()))
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

	t.Run("generate token and parse token for moderator", func(t *testing.T) {
		ctx := context.Background()
		user := &entity.User{
			Email:    "test",
			Password: "test",
			UserType: "moderator",
		}

		log.Debug("creating user")
		userID, err := authService.CreateUser(ctx, &AuthCreateUserInput{
			Email:    user.Email,
			Password: user.Password,
			UserType: user.UserType,
		})
		require.NoError(t, err)
		require.True(t, len(userID) > 0)
		user.ID = userID

		log.Debug("generating token")
		generatedToken, err := authService.GenerateToken(ctx, &AuthGenerateTokenInput{
			Email:    user.Email,
			Password: user.Password,
		})
		require.NoError(t, err)

		log.Debug("parsing token")
		claims, err := authService.ParseToken(generatedToken)
		require.NoError(t, err)
		require.Equal(t, claims.UserType, user.UserType)
		require.Equal(t, claims.UserType, "moderator")
		require.Equal(t, claims.UserID, user.ID)
	})

	t.Run("generate token for user with empty email", func(t *testing.T) {
		ctx := context.Background()
		user := &entity.User{
			Password: "test2",
			UserType: "moderator",
		}

		log.Debug("generating token")
		generatedToken, err := authService.GenerateToken(ctx, &AuthGenerateTokenInput{
			Email:    user.Email,
			Password: user.Password,
		})
		require.ErrorIs(t, err, ErrUserNotFound)
		require.Equal(t, "", generatedToken)
	})

	t.Run("generate token for user with invalid password", func(t *testing.T) {
		ctx := context.Background()
		user := &entity.User{
			Email:    "test",
			Password: "test2",
			UserType: "moderator",
		}

		log.Debug("generating token")
		generatedToken, err := authService.GenerateToken(ctx, &AuthGenerateTokenInput{
			Email:    user.Email,
			Password: user.Password,
		})
		require.ErrorIs(t, err, ErrWrongPassword)
		require.Equal(t, "", generatedToken)
	})

}
