package domain

import (
	"context"
)

type UserProfile struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id" validate:"required"`
	NIK          string `json:"nik" validate:"required"`
	FullName     string `json:"full_name" validate:"required"`
	LegalName    string `json:"legal_name" validate:"required"`
	PlaceOfBirth string `json:"pob" validate:"required"`
	DateOfBirth  string `json:"dob" validate:"required"`
	Salary       int    `json:"salary" validate:"required"`
	KTP          string `json:"ktp" validate:"required"`
	Selfie       string `json:"selfie" validate:"required"`
}

type ServiceUserProfile interface {
	// Adding service
	// AddUserProfile is a method to add a new user_profile to the database
	AddUserProfile(ctx context.Context, up *UserProfile) (int, error)

	// Listing service
	// GetUserProfileByUserID is a method to get user_profile by user_id from the database
	GetUserProfileByUserID(ctx context.Context, userID int) (UserProfile, error)
}

type RepositoryMySQLUserProfile interface {
	// Adding repository
	// AddUserProfile is a method to add a new user_profile to the database
	AddUserProfile(ctx context.Context, up *UserProfile) (int, error)

	// Listing repository
	// ReadUserProfileByUserID is a method to get user_profile by user_id from the database
	ReadUserProfileByUserID(ctx context.Context, userID int) (UserProfile, error)
}
