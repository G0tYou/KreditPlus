package service

import (
	d "app/domain"
	"context"
)

type serviceLimitType struct {
	rmsqllt d.RepositoryMySQLLimitType
}

// NewService creates an service with the necessary dependencies
func NewServiceLimitType(rmsqllt d.RepositoryMySQLLimitType) d.ServiceLimitType {
	return &serviceLimitType{rmsqllt}
}

// Write the adding service below
// AddLimitType is a method to add a new limit_type to the reporsitory
func (s *serviceLimitType) AddLimitType(ctx context.Context, lt *d.LimitType) (int, error) {
	id, err := s.rmsqllt.AddLimitType(ctx, lt)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Write the listing service below
// GetLimitTypes is a method to get list of limit_type from the repository
func (s *serviceLimitType) GetLimitTypes(ctx context.Context) ([]d.LimitType, error) {
	lts, err := s.rmsqllt.ReadLimitTypes(ctx)
	if err != nil {
		return nil, err
	}
	return lts, nil
}
