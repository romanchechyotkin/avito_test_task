package request

type Login struct {
	Email    string `json:"email" validate:"email" example:"test@gmail.com"`
	Password string `json:"password" validate:"min=4,max=50" example:"password"`
}
