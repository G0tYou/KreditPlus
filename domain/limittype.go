package domain

import "context"

type LimitType struct {
	ID    int `json:"id"`
	Tenor int `json:"tenor"`
}

type ServiceLimitType interface {
	AddLimitType(ctx context.Context, lt *LimitType) (int, error)
	ReadLimitTypes(ctx context.Context) ([]LimitType, error)
}

type RepositoryMySQLLimitType interface {
	AddLimitType(ctx context.Context, lt *LimitType) (int, error)
	ReadLimitTypes(ctx context.Context) ([]LimitType, error)
}
