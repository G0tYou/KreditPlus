package repository

import (
	d "app/domain"
	"context"
	"database/sql"
)

type mysqlUserProfileRepository struct {
	Conn *sql.DB
}

// NewMysqlLimitTypeRepository creates a new instance of mysqlLimitTypeRepository
func NewMysqlRepositoryUserProfile(conn *sql.DB) d.RepositoryMySQLUserProfile {
	return &mysqlUserProfileRepository{conn}
}

// Write the adding repository below
// AddUserProfile is a method to add a new user_profile to the database
func (m *mysqlUserProfileRepository) AddUserProfile(ctx context.Context, up *d.UserProfile) (int, error) {
	query := "INSERT INTO user_profile (user_id, nik, full_name, legal_name, place_of_birth, date_of_birth, salary, ktp, selfie) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, up.UserID, up.NIK, up.FullName, up.LegalName, up.PlaceOfBirth, up.DateOfBirth, up.Salary, up.KTP, up.Selfie)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}
