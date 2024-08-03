package flat

import (
	"context"
	"log/slog"
	"time"

	"github.com/romanchechyotkin/avito_test_task/internal/entity"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"
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
	var err error

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		r.log.Debug("failed to start transaction", logger.Error(err))
		return nil, err
	}

	defer func() {
		if err != nil {
			r.log.Debug("rollbacking transaction")
			if err = tx.Rollback(ctx); err != nil {
				r.log.Error("failed to rollback transaction", logger.Error(err))
				return
			}
		} else {
			r.log.Debug("committing transaction")
			if err = tx.Commit(ctx); err != nil {
				r.log.Error("failed to commit transaction", logger.Error(err))
				return
			}
		}
	}()

	q := `INSERT INTO flats (number, house_id, price, rooms_amount) VALUES ($1, $2, $3, $4)
	RETURNING id, number, house_id, price, rooms_amount, moderation_status, created_at
`

	r.log.Debug("create flat query", slog.String("query", q))

	if err = tx.QueryRow(ctx, q, flat.Number, flat.HouseID, flat.Price, flat.RoomsAmount).Scan(
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

	q = `UPDATE houses SET updated_at = $1 WHERE id = $2`

	r.log.Debug("update house query", slog.String("query", q))
	exec, err := tx.Exec(ctx, q, time.Now(), flat.HouseID)
	if err != nil {
		return nil, err
	}

	r.log.Debug("update result", slog.Int64("rows affected", exec.RowsAffected()))

	return flat, nil
}
