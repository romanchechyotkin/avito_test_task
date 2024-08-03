package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/romanchechyotkin/avito_test_task/internal/repo"
)

type AuthCreateUserInput struct {
	Email    string
	Password string
	UserType string
}

type AuthGenerateTokenInput struct {
	Email    string
	Password string
}

type Auth interface {
	CreateUser(ctx context.Context, input *AuthCreateUserInput) (int, error)

	GenerateToken(ctx context.Context, input *AuthGenerateTokenInput) (string, error)

	ParseToken(accessToken string) (*TokenClaims, error)
}

type Dependencies struct {
	Log   *slog.Logger
	Repos *repo.Repositories

	//Hasher hasher.PasswordHasher
	//
	SignKey  string
	TokenTTL time.Duration
}

type Services struct {
	Auth Auth
}

func NewServices(deps *Dependencies) *Services {
	return &Services{
		Auth: NewAuthService(deps.Log, deps.Repos.User, deps.SignKey, deps.TokenTTL),
	}
}
