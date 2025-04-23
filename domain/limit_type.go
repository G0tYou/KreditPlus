package domain

import "context"

type LimitType struct {
	ID    int `json:"id"`
	Tenor int `json:"tenor" validate:"required"`
}

type ServiceLimitType interface {
	// Adding service
	// AddLimitType is a method to add a new limit type to the database
	AddLimitType(ctx context.Context, lt *LimitType) (int, error)
	// Listing service
	// GetLimitTypes is a method to get all limit types from the database
	GetLimitTypes(ctx context.Context) ([]LimitType, error)
}

type RepositoryMySQLLimitType interface {
	// Adding repository
	// AddLimitType is a method to add a new limit type to the database
	AddLimitType(ctx context.Context, lt *LimitType) (int, error)

	// Listing repository
	// ReadLimitTypes is a method to get all limit types from the database
	ReadLimitTypes(ctx context.Context) ([]LimitType, error)
}
