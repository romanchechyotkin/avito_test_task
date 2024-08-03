package request

type CreateFlat struct {
	Number      uint `json:"number" validate:"required"`
	HouseID     uint `json:"house_id" validate:"required"`
	Price       uint `json:"price" validate:"required"`
	RoomsAmount uint `json:"rooms_amount" validate:"required"`
}
