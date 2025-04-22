package domain

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ServiceUser interface {
}

type RepositoryUser interface {
}
