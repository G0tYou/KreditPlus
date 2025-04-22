package service

import (
	d "app/domain"
	h "app/internal/helper"
	"context"
)

type serviceUser struct {
	rmsqlu d.RepositoryMySQLUser
}

// NewService creates an service with the necessary dependencies
func NewServiceUser(rmsqlu d.RepositoryMySQLUser) d.ServiceUser {
	return &serviceUser{rmsqlu}
}

// Write the adding service below
// AddUser is a method to add a new user to the repository
func (s *serviceUser) AddUser(ctx context.Context, u *d.User) (int, error) {
	var err error

	//validate username is exist
	exist, err := s.rmsqlu.ExistByUsername(ctx, u.Username)
	if err != nil {
		return 0, err
	}

	if !exist {
		//encrypt password
		u.Password, err = h.EncodePassword(u.Password)
		if err != nil {
			return 0, err
		}

		id, err := s.rmsqlu.AddUser(ctx, u)
		if err != nil {
			return 0, err
		}

		return id, nil
	}

	return 0, d.ErrConflictUsername
}
