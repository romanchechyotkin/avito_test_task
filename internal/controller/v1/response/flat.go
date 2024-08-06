package response

import (
	"github.com/romanchechyotkin/avito_test_task/internal/entity"
	"time"
)

type HouseFlats struct {
	Flats []*entity.Flat `json:"flats"`
}

type Flat struct {
	ID               uint      `json:"id"`
	Number           uint      `json:"number"`
	HouseID          uint      `json:"house_id"`
	Price            uint      `json:"price"`
	RoomsAmount      uint      `json:"rooms_amount"`
	ModerationStatus string    `json:"moderation_status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func BuildFlat(flat *entity.Flat) Flat {
	return Flat{
		ID:               flat.ID,
		Number:           flat.Number,
		HouseID:          flat.HouseID,
		Price:            flat.Price,
		RoomsAmount:      flat.RoomsAmount,
		ModerationStatus: flat.ModerationStatus,
		CreatedAt:        flat.CreatedAt,
		UpdatedAt:        flat.UpdatedAt,
	}
}
