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
	// Adding service
	// AddUser is a service to add a new user to the repository
	AddUser(ctx context.Context, u *User) (int, error)

	//validating service
	// Login is a service to validate data user from the repository
	Login(ctx context.Context, u *User) (string, error)
}

type RepositoryMySQLUser interface {
	// Adding repository
	// AddUser is a repository to add a new user to the database
	AddUser(ctx context.Context, u *User) (int, error)

	// Listing repository
	// ReadUserByUsername is a repository to get data user from the database
	ReadUserByUsername(ctx context.Context, username string) (*User, error)

	//validating repository below
	//ExistByUsername is a repository to validate user is exist in the database
	ExistByUsername(ctx context.Context, username string) (bool, error)
}
