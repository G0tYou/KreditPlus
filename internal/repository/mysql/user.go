package repository

import (
	"context"
	"database/sql"

	d "app/domain"
)

type mysqlRepositoryUser struct {
	Conn *sql.DB
}

// NewMysqlRepositoryUser creates a new instance of mysqlLimitTypeRepository
func NewMysqlRepositoryUser(conn *sql.DB) d.RepositoryMySQLUser {
	return &mysqlRepositoryUser{conn}
}

// Write the adding repository below
// AddUser adds a new user to the database
func (m *mysqlRepositoryUser) AddUser(ctx context.Context, u *d.User) (int, error) {
	query := "INSERT user SET username = ?, password = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, u.Username, u.Password)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Write the validating repository below
// ExistByUsername is a method to validate user is exist in the database
func (m *mysqlRepositoryUser) ExistByUsername(ctx context.Context, username string) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM user WHERE username = ?)"

	err := m.Conn.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, err
}
