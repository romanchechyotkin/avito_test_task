package request

type CreateFlat struct {
	Number      uint `json:"number" validate:"required" example:"test@gmail.com"`
	HouseID     uint `json:"house_id" validate:"required" example:"123"`
	Price       uint `json:"price" validate:"required" example:"22344"`
	RoomsAmount uint `json:"rooms_amount" validate:"required" example:"3"`
}

type UpdateFlat struct {
	ID     uint   `json:"id" validate:"required" example:"123"`
	Status string `json:"status" validate:"required,oneof=created approved declined 'on moderation'" enums:"created,approved,declined,on moderation" example:"on moderation"`
}
