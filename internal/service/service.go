package service

type Auth interface {
}

type Services struct {
	Auth Auth
}

func NewServices() *Services {
	return &Services{
		Auth: NewAuthService(),
	}
}
