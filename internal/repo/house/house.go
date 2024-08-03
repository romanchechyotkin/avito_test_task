package house

import (
	"context"
	"log/slog"

	"github.com/romanchechyotkin/avito_test_task/internal/entity"
	"github.com/romanchechyotkin/avito_test_task/pkg/postgresql"
)

type Repo struct {
	log *slog.Logger
	*postgresql.Postgres
}

func NewRepo(log *slog.Logger, pg *postgresql.Postgres) *Repo {
	return &Repo{
		log:      log,
		Postgres: pg,
	}
}

func (r *Repo) CreateHouse(ctx context.Context, house *entity.House) (*entity.House, error) {
	q := `INSERT INTO houses (address, year, developer) VALUES ($1, $2, $3)
	RETURNING id, address, year, developer, created_at, updated_at
`

	r.log.Debug("create house query", slog.String("query", q))

	if err := r.Pool.QueryRow(ctx, q, house.Address, house.Year, house.Developer).Scan(
		&house.ID,
		&house.Address,
		&house.Year,
		&house.Developer,
		&house.CreatedAt,
		&house.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return house, nil
}
