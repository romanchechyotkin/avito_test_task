package repo

import (
	"context"
	"github.com/romanchechyotkin/avito_test_task/internal/repo/house"
	"log/slog"

	"github.com/romanchechyotkin/avito_test_task/internal/entity"
	"github.com/romanchechyotkin/avito_test_task/internal/repo/user"
	"github.com/romanchechyotkin/avito_test_task/pkg/postgresql"
)

type User interface {
	CreateUser(ctx context.Context, user *entity.User) (string, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetById(ctx context.Context, id int) (*entity.User, error)
}

type House interface {
	CreateHouse(ctx context.Context, house *entity.House) (*entity.House, error)
}

type Repositories struct {
	User
	House
}

func NewRepositories(log *slog.Logger, pg *postgresql.Postgres) *Repositories {
	return &Repositories{
		User:  user.NewRepo(log, pg),
		House: house.NewRepo(log, pg),
	}
}
