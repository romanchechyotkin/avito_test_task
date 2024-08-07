package response

type Login struct {
	Token string `json:"token" example:"auth token"`
}
