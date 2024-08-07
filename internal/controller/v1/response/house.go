package response

import (
	"time"

	"github.com/romanchechyotkin/avito_test_task/internal/entity"
)

func BuildHouse(house *entity.House) House {
	var developer string

	if house.Developer.Valid {
		developer = house.Developer.String
	}

	return House{
		ID:        house.ID,
		Address:   house.Address,
		Year:      house.Year,
		Developer: developer,
		CreatedAt: house.CreatedAt,
		UpdatedAt: house.UpdatedAt,
	}
}

type House struct {
	ID        uint      `json:"id" example:"123"`
	Address   string    `json:"address" example:"Улица Пушкина 1"`
	Year      uint      `json:"year" example:"1999"`
	Developer string    `json:"developer,omitempty" example:"ООО Компания"`
	CreatedAt time.Time `json:"created_at" example:"2024-08-09T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-08-09T00:00:00Z"`
}
