package domain

import (
	"context"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ServiceUser interface {
	AddUser(ctx context.Context, u *User) (int, error)
}

type RepositoryMySQLUser interface {
	//adding repository below
	AddUser(ctx context.Context, u *User) (int, error)
	//validating repository below
	ExistByUsername(ctx context.Context, username string) (bool, error)
}
