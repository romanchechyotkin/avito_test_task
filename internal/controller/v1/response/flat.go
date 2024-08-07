package response

import (
	"time"

	"github.com/romanchechyotkin/avito_test_task/internal/entity"
)

type HouseFlats struct {
	Flats []*entity.Flat `json:"flats"`
}

type Flat struct {
	ID               uint      `json:"id" example:"123"`
	Number           uint      `json:"number" example:"1"`
	HouseID          uint      `json:"house_id" example:"1"`
	Price            uint      `json:"price" example:"123"`
	RoomsAmount      uint      `json:"rooms_amount" example:"4"`
	ModerationStatus string    `json:"moderation_status" example:"created"`
	CreatedAt        time.Time `json:"created_at" example:"2024-08-09T00:00:00Z"`
	UpdatedAt        time.Time `json:"updated_at" example:"2024-08-09T00:00:00Z"`
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
