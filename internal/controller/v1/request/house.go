package request

type CreateHouse struct {
	Address   string `json:"address" validate:"required" example:"ул. Новая, д. 1"`
	Year      uint   `json:"year" validate:"required" example:"2022"`
	Developer string `json:"developer" example:"ООО Компания"`
}
