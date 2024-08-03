package entity

import (
	"database/sql"
	"time"
)

type House struct {
	ID        uint           `db:"id"`
	Address   string         `db:"address"`
	Year      uint           `db:"year"`
	Developer sql.NullString `db:"developer"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}
