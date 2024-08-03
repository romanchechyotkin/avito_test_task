package entity

import (
	"time"
)

type Flat struct {
	ID               uint      `json:"id" db:"id"`
	Number           uint      `json:"number" db:"number"`
	HouseID          uint      `json:"house_id" db:"address"`
	Price            uint      `json:"price" db:"price"`
	RoomsAmount      uint      `json:"rooms_amount" db:"rooms_amount"`
	ModerationStatus string    `json:"moderation_status" db:"moderation_status"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}
