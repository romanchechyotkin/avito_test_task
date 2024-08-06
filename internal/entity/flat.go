package entity

import (
	"time"
)

type Flat struct {
	ID               uint      `db:"id"`
	Number           uint      `db:"number"`
	HouseID          uint      `db:"address"`
	Price            uint      `db:"price"`
	RoomsAmount      uint      `db:"rooms_amount"`
	ModerationStatus string    `db:"moderation_status"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
