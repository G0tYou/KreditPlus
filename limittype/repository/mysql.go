package repository

import (
	"context"
	"database/sql"

	d "app/domain"
)

type mysqlLimitTypeRepository struct {
	Conn *sql.DB
}

// NewMysqlLimitTypeRepository creates a new instance of mysqlLimitTypeRepository
func NewMysqlLimitTypeRepository(conn *sql.DB) d.RepositoryMySQLLimitType {
	return &mysqlLimitTypeRepository{conn}
}

// Write the adding repository below
// AddLimitType is a method to add a new limit_type to the database
func (m *mysqlLimitTypeRepository) AddLimitType(ctx context.Context, lt *d.LimitType) (int, error) {
	query := "INSERT limit_type SET tenor = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, lt.Tenor)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Write the listing repository below
// ReadLimitTypes is a method to get list of limit_type from the database
func (m *mysqlLimitTypeRepository) ReadLimitTypes(ctx context.Context) ([]d.LimitType, error) {
	query := "SELECT id, tenor FROM limit_type"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lts []d.LimitType
	for rows.Next() {
		var lt d.LimitType
		if err := rows.Scan(&lt.ID, &lt.Tenor); err != nil {
			return nil, err
		}
		lts = append(lts, lt)
	}

	return lts, nil
}
