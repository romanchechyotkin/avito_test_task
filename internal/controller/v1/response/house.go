package response

import (
	"time"

	"github.com/romanchechyotkin/avito_test_task/internal/entity"
)

func BuildHouse(house *entity.House) House {
	var developer string
	var updatedAt time.Time

	if house.Developer.Valid {
		developer = house.Developer.String
	}

	return House{
		ID:        house.ID,
		Address:   house.Address,
		Year:      house.Year,
		Developer: developer,
		CreatedAt: house.CreatedAt,
		UpdatedAt: updatedAt,
	}
}

type House struct {
	ID        uint      `json:"id"`
	Address   string    `json:"address"`
	Year      uint      `json:"year"`
	Developer string    `json:"developer,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_ats"`
}
