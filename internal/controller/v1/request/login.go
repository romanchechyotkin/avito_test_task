package request

type Login struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=4,max=50"`
}
