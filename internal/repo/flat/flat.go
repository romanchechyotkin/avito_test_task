package flat

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

func (r *Repo) CreateFlat(ctx context.Context, flat *entity.Flat) (*entity.Flat, error) {
	q := `INSERT INTO flats (number, house_id, price, rooms_amount) VALUES ($1, $2, $3, $4)
	RETURNING id, number, house_id, price, rooms_amount, moderation_status, created_at
`

	r.log.Debug("create flat query", slog.String("query", q))

	if err := r.Pool.QueryRow(ctx, q, flat.Number, flat.HouseID, flat.Price, flat.RoomsAmount).Scan(
		&flat.ID,
		&flat.Number,
		&flat.HouseID,
		&flat.Price,
		&flat.RoomsAmount,
		&flat.ModerationStatus,
		&flat.CreatedAt,
	); err != nil {
		return nil, err
	}

	return flat, nil
}
