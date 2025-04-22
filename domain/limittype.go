package domain

type LimitType struct {
	ID    int `json:"id"`
	Tenor int `json:"tenor"`
}

type Service interface{}

type RepositoryMySQL interface{}
