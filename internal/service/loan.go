package service

import (
	d "app/domain"
	"context"
)

type serviceLoan struct {
	rmsqllo d.RepositoryMySQLLoan
	rmsqli  d.RepositoryMySQLLimit
}

// NewService creates an service with the necessary dependencies
func NewServiceLoan(rmsqllo d.RepositoryMySQLLoan, rmsqli d.RepositoryMySQLLimit) d.ServiceLoan {
	return &serviceLoan{rmsqllo, rmsqli}
}

// AddLoan is a method to add a new loan to the repository
func (s *serviceLoan) AddLoan(ctx context.Context, lo *d.Loan) (int, error) {

	tx, err := s.rmsqllo.BeginTx(ctx)
	if err != nil {
		return 0, d.ErrConflictUsername
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	li, err := s.rmsqli.ReadTXLimitByID(ctx, tx, lo.LimitID)
	if err != nil {
		return 0, err
	}

	if li.UsedAmount+lo.Amount > li.Amount {
		return 0, d.ErrLimitExceeded
	}

	id, err := s.rmsqllo.AddTXLoan(ctx, tx, lo)
	if err != nil {
		return 0, err
	}

	err = s.rmsqli.UpdateTXLimit(ctx, tx, &d.Limit{ID: li.ID, UsedAmount: li.UsedAmount + lo.Amount})
	if err != nil {
		return 0, err
	}

	return id, nil
}
