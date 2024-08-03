package repo

import (
	"context"
	"log/slog"

	"github.com/romanchechyotkin/avito_test_task/internal/entity"
	"github.com/romanchechyotkin/avito_test_task/internal/repo/user"
	"github.com/romanchechyotkin/avito_test_task/pkg/postgresql"
)

type User interface {
	CreateUser(ctx context.Context, user *entity.User) (int, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetById(ctx context.Context, id int) (*entity.User, error)
}

type Repositories struct {
	User
}

func NewRepositories(log *slog.Logger, pg *postgresql.Postgres) *Repositories {
	return &Repositories{
		User: user.NewRepo(log, pg),
	}
}
