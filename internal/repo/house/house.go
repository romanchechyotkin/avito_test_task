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

func (r *Repo) CreateSubscription(ctx context.Context, houseID, userID string) error {
	q := "INSERT INTO house_subscriptions (house_id, user_id) VALUES ($1, $2)"

	r.log.Debug("create house subscription query", slog.String("query", q))

	exec, err := r.Pool.Exec(ctx, q, houseID, userID)
	if err != nil {
		// todo 23503 violates foreign key constraint
		// todo violates unique constraint \"house_subscriptions_pkey\" (SQLSTATE 23505

		return err
	}

	r.log.Debug("exec result", slog.Int64("rows affected", exec.RowsAffected()))

	return nil
}

func (r *Repo) GetHouseSubscriptions(ctx context.Context, houseID uint) ([]string, error) {
	q := "SELECT u.email FROM users u JOIN house_subscriptions hs ON u.id = hs.user_id WHERE hs.house_id = $1"

	rows, err := r.Pool.Query(ctx, q, houseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []string

	for rows.Next() {
		var email string

		if err = rows.Scan(&email); err != nil {
			return nil, err
		}

		emails = append(emails, email)
	}

	return emails, nil
}
