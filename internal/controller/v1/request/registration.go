package request

type Registration struct {
	Email    string `json:"email" validate:"email" example:"test@gmail.com"`
	Password string `json:"password" validate:"min=4,max=50" example:"password"`
	UserType string `json:"user_type" validate:"oneof=client moderator" enums:"client,moderator" example:"client"`
}
