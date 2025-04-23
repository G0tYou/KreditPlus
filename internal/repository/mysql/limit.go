package repository

import (
	"context"
	"database/sql"

	d "app/domain"
)

type mysqlRepositoryLimit struct {
	Conn *sql.DB
}

// NewMysqlLimitTypeRepository creates a new instance of mysqlLimitTypeRepository
func NewMysqlRepositoryLimit(conn *sql.DB) d.RepositoryMySQLLimit {
	return &mysqlRepositoryLimit{conn}
}

// Write the listing repository below
// ReadLimitByUserIDAndLimitTypeID is a method to get limit by user_id and limit_type_id from the database
func (m *mysqlRepositoryLimit) ReadTXLimitByID(ctx context.Context, tx *sql.Tx, lid int) (d.Limit, error) {
	var l d.Limit

	query := "SELECT id, limit_type_id, user_id, amount, used_amount FROM `limit` WHERE id = ? FOR UPDATE"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return l, err
	}
	defer stmt.Close()

	err = tx.QueryRowContext(ctx, query, lid).Scan(&l.ID, &l.LimitTypeID, &l.UserID, &l.Amount, &l.UsedAmount)
	if err != nil {
		return l, err
	}

	return l, nil
}

// Write the updating repository below
// UpdateLimit is a method to update limit in the database
func (m *mysqlRepositoryLimit) UpdateTXLimit(ctx context.Context, tx *sql.Tx, l *d.Limit) error {
	query := "UPDATE `limit` SET used_amount = ? WHERE id = ?"

	_, err := tx.ExecContext(ctx, query, l.UsedAmount, l.ID)
	if err != nil {
		return err
	}

	return nil
}
