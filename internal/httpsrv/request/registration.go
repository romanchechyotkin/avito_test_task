package request

type Registration struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=4,max=50"`
	UserType string `json:"user_type" validate:"oneof=client moderator"`
}
