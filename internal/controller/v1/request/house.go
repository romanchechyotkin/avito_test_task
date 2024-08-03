package request

type CreateHouse struct {
	Address   string `json:"address" validate:"required"`
	Year      uint   `json:"year" validate:"required"`
	Developer string `json:"developer"`
}
