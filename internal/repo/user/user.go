package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/romanchechyotkin/avito_test_task/internal/entity"
	"github.com/romanchechyotkin/avito_test_task/internal/repo/codes"
	"github.com/romanchechyotkin/avito_test_task/internal/repo/repoerrors"
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

func (r *Repo) CreateUser(ctx context.Context, user *entity.User) (string, error) {
	q := "INSERT INTO users (email, password, user_type) VALUES ($1, $2, $3) RETURNING id"

	r.log.Debug("create user query", slog.String("query", q))

	var id string
	if err := r.Pool.QueryRow(ctx, q, user.Email, user.Password, user.UserType).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == codes.UniqueConstraintCode {
				return "", repoerrors.ErrAlreadyExists
			}
		}

		return "", err
	}

	return id, nil
}

func (r *Repo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	q := "SELECT id, email, password, user_type, created_at FROM users WHERE email = $1"

	r.log.Debug("get user by email query", slog.String("query", q))

	var user entity.User
	err := r.Pool.QueryRow(ctx, q, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.UserType,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrors.ErrNotFound
		}

		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == codes.UniqueConstraintCode {
				return nil, repoerrors.ErrAlreadyExists
			}
		}

		return nil, err
	}

	return &user, nil
}

func (r *Repo) GetById(ctx context.Context, id int) (*entity.User, error) {
	q := "SELECT id, email, password, user_type, created_at FROM users WHERE id = $1"

	r.log.Debug("get user by id query", slog.String("query", q))

	var user entity.User
	err := r.Pool.QueryRow(ctx, q, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.UserType,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrors.ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}
