package domain

import (
	"context"
	"database/sql"
)

type Loan struct {
	ID      int `json:"id"`
	LimitID int `json:"limit_id" validate:"required"`
	Amount  int `json:"amount" validate:"required"`
}

type ServiceLoan interface {
	//Adding service
	// AddLoan is a method to add a new loan to the repository
	AddLoan(ctx context.Context, lo *Loan) (int, error)

	//Listing service
}

type RepositoryMySQLLoan interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)

	//Adding repository
	// AddTXLoan is a method to add loan in the database
	AddTXLoan(ctx context.Context, tx *sql.Tx, l *Loan) (int, error)

	//Listing repository
}
