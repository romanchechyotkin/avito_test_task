package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/romanchechyotkin/avito_test_task/internal/entity"
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
	CreateUser(ctx context.Context, input *AuthCreateUserInput) (string, error)

	GenerateToken(ctx context.Context, input *AuthGenerateTokenInput) (string, error)

	ParseToken(accessToken string) (*TokenClaims, error)
}

type HouseCreateInput struct {
	Address   string
	Year      uint
	Developer string
}

type House interface {
	CreateHouse(ctx context.Context, input *HouseCreateInput) (*entity.House, error)
}

type Dependencies struct {
	Log   *slog.Logger
	Repos *repo.Repositories

	SignKey  string
	TokenTTL time.Duration
}

type Services struct {
	Auth  Auth
	House House
}

func NewServices(deps *Dependencies) *Services {
	return &Services{
		Auth:  NewAuthService(deps.Log, deps.Repos.User, deps.SignKey, deps.TokenTTL),
		House: NewHouseService(deps.Log, deps.Repos.House),
	}
}
