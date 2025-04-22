package domain

type LimitType struct {
	ID    int `json:"id"`
	Tenor int `json:"tenor"`
}

type ServiceLimitType interface{}

type RepositoryMySQLLimitType interface{}
