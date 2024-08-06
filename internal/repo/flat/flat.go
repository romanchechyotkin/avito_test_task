package flat

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/romanchechyotkin/avito_test_task/internal/entity"
	"github.com/romanchechyotkin/avito_test_task/internal/repo/codes"
	"github.com/romanchechyotkin/avito_test_task/internal/repo/repoerrors"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"
	"github.com/romanchechyotkin/avito_test_task/pkg/postgresql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	RETURNING id, number, house_id, price, rooms_amount, moderation_status, created_at, updated_at
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
		&flat.UpdatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == codes.ForeignKeyConstraint {
				return nil, repoerrors.ErrNotFound
			}
		}

		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == codes.UniqueConstraintCode {
				return nil, repoerrors.ErrAlreadyExists
			}
		}

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

func (r *Repo) GetStatus(ctx context.Context, id uint) (string, error) {
	q := `SELECT moderation_status FROM flats WHERE id = $1`

	r.log.Debug("select flat status query", slog.String("query", q))

	var status string
	if err := r.Pool.QueryRow(ctx, q, id).Scan(&status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repoerrors.ErrNotFound
		}

		return "", err
	}

	return status, nil
}

func (r *Repo) UpdateStatus(ctx context.Context, flat *entity.Flat, moderatorID sql.NullString) (*entity.Flat, error) {
	q := `UPDATE flats SET moderation_status = $1, moderator_id = $2, updated_at = $3 WHERE id = $4
	RETURNING id, number, house_id, price, rooms_amount, moderation_status, created_at, updated_at
`

	r.log.Debug("update flat status query", slog.String("query", q))

	if err := r.Pool.QueryRow(ctx, q, flat.ModerationStatus, moderatorID, time.Now(), flat.ID).Scan(
		&flat.ID,
		&flat.Number,
		&flat.HouseID,
		&flat.Price,
		&flat.RoomsAmount,
		&flat.ModerationStatus,
		&flat.CreatedAt,
		&flat.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrors.ErrNotFound
		}

		return nil, err
	}

	return flat, nil
}

// todo use lib for binding values
func (r *Repo) GetHouseFlats(ctx context.Context, houseID, userType string) ([]*entity.Flat, error) {
	q := `SELECT id, number, house_id, price, rooms_amount, moderation_status, created_at, updated_at FROM flats WHERE house_id = $1 `

	if userType == "client" {
		q += "AND moderation_status = 'approved'"
	}

	r.log.Debug("select all flats for house query", slog.String("query", q))

	rows, err := r.Pool.Query(ctx, q, houseID)
	if err != nil {
		r.log.Debug("failed to query rows", logger.Error(err))
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrors.ErrNotFound
		}

		return nil, err
	}
	defer rows.Close()

	var flats []*entity.Flat

	for rows.Next() {
		var flat entity.Flat

		err = rows.Scan(
			&flat.ID,
			&flat.Number,
			&flat.HouseID,
			&flat.Price,
			&flat.RoomsAmount,
			&flat.ModerationStatus,
			&flat.CreatedAt,
			&flat.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		flats = append(flats, &flat)
	}

	return flats, nil
}
