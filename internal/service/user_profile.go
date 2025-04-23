package service

import (
	d "app/domain"
	"context"
	"time"
)

type serviceUserProfile struct {
	rmsqlup d.RepositoryMySQLUserProfile
}

// NewService creates an service with the necessary dependencies
func NewServiceUserProfile(rmsqlup d.RepositoryMySQLUserProfile) d.ServiceUserProfile {
	return &serviceUserProfile{rmsqlup}
}

// Write the adding service below
// AddUserProfile is a method to add a new user_profile to the repository
func (s *serviceUserProfile) AddUserProfile(ctx context.Context, up *d.UserProfile) (int, error) {
	birtdate, _ := time.Parse("2006/01/02", up.DateOfBirth)

	// Format the date to YYYY-MM-DD
	up.DateOfBirth = birtdate.Format("2006-01-02")

	id, err := s.rmsqlup.AddUserProfile(ctx, up)
	if err != nil {
		return 0, err
	}

	return id, nil
}
