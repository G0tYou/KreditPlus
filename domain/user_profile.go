package domain

import (
	"context"
	"database/sql"
)

type UserProfile struct {
	ID           int    `json:"id"`
	UserID       string `json:"user_id" validate:"required"`
	NIK          string `json:"nik" validate:"required"`
	FullName     string `json:"full_name" validate:"required"`
	LegalName    string `json:"legal_name" validate:"required"`
	PlaceOfBirth string `json:"pob" validate:"required"`
	DateOfBirth  string `json:"dob" validate:"required"`
	Salary       string `json:"alary" validate:"required"`
	KTP          string `json:"ktp" validate:"required"`
	Selfie       string `json:"selfie" validate:"required"`
}

type RepositoryMySQLUserProfile interface {
	AddUserProfile(ctx context.Context, up *UserProfile, tx *sql.Tx) (int, error)
}
