package repository

import (
	"database/sql"

	d "app/domain"
)

type mysqlUserRepository struct {
	Conn *sql.DB
}

// NewMysqlLimitTypeRepository creates a new instance of mysqlLimitTypeRepository
func NewMysqlUserRepository(conn *sql.DB) d.RepositoryUser {
	return &mysqlUserRepository{conn}
}
