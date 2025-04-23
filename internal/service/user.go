package service

import (
	d "app/domain"
	h "app/internal/helper"
	"context"

	"github.com/spf13/viper"
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

// Write the validating service below
// Login is a method to get data user from the repository
func (s *serviceUser) Login(ctx context.Context, u *d.User) (string, error) {
	//validate username is exist
	exist, err := s.rmsqlu.ExistByUsername(ctx, u.Username)
	if err != nil {
		return "", err
	}

	if exist {
		lu, err := s.rmsqlu.ReadUserByUsername(ctx, u.Username)
		if err != nil {
			return "", err
		}

		//compare password
		err = h.ComparePassword(lu.Password, u.Password)
		if err != nil {
			return "", d.ErrLogin
		}

		//generate token
		token, err := h.GenerateJWT(int64(lu.ID), viper.GetString("jwt.secret"))
		if err != nil {
			return "", err
		}

		return token, nil
	}

	return "", d.ErrUserNotFound
}
