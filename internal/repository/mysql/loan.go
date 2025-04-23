package repository

import (
	"context"
	"database/sql"

	d "app/domain"
)

type mysqlRepositoryLoan struct {
	Conn *sql.DB
}

// NewMysqlLimitTypeRepository creates a new instance of mysqlLimitTypeRepository
func NewMysqlRepositoryLoan(conn *sql.DB) d.RepositoryMySQLLoan {
	return &mysqlRepositoryLoan{conn}
}

func (r *mysqlRepositoryLoan) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.Conn.BeginTx(ctx, nil)
}

// Write the Adding repository below
// AddTXLoan is a method to add loan in the database
func (r *mysqlRepositoryLoan) AddTXLoan(ctx context.Context, tx *sql.Tx, l *d.Loan) (int, error) {
	query := "INSERT INTO loan (limit_id, amount) VALUES(?,?)"

	res, err := tx.ExecContext(ctx, query, l.LimitID, l.Amount)
	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()

	return int(id), nil
}
