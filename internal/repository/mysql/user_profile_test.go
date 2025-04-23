package repository_test

import (
	"app/domain"
	mysql "app/internal/repository/mysql"
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddUserProfile(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Create a new repository instance
	repo := mysql.NewMysqlRepositoryUserProfile(db)

	// Define test cases
	tests := []struct {
		name          string
		input         *domain.UserProfile
		mockSetup     func()
		expectedID    int
		expectedError error
	}{
		{
			name: "Success",
			input: &domain.UserProfile{
				UserID:       1,
				NIK:          "123456789",
				FullName:     "John Doe",
				LegalName:    "Johnathan Doe",
				PlaceOfBirth: "City A",
				DateOfBirth:  "1990-01-01",
				Salary:       5000,
				KTP:          "ktp123",
				Selfie:       "selfie123",
			},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT INTO user_profile").
					ExpectExec().
					WithArgs(1, "123456789", "John Doe", "Johnathan Doe", "City A", "1990-01-01", 5000, "ktp123", "selfie123").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedID:    1,
			expectedError: nil,
		},
		{
			name: "PrepareContext Error",
			input: &domain.UserProfile{
				UserID:       1,
				NIK:          "123456789",
				FullName:     "John Doe",
				LegalName:    "Johnathan Doe",
				PlaceOfBirth: "City A",
				DateOfBirth:  "1990-01-01",
				Salary:       5000,
				KTP:          "ktp123",
				Selfie:       "selfie123",
			},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT INTO user_profile").
					WillReturnError(sql.ErrConnDone)
			},
			expectedID:    0,
			expectedError: sql.ErrConnDone,
		},
		{
			name: "ExecContext Error",
			input: &domain.UserProfile{
				UserID:       1,
				NIK:          "123456789",
				FullName:     "John Doe",
				LegalName:    "Johnathan Doe",
				PlaceOfBirth: "City A",
				DateOfBirth:  "1990-01-01",
				Salary:       5000,
				KTP:          "ktp123",
				Selfie:       "selfie123",
			},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT INTO user_profile").
					ExpectExec().
					WithArgs(1, "123456789", "John Doe", "Johnathan Doe", "City A", "1990-01-01", 5000, "ktp123", "selfie123").
					WillReturnError(sql.ErrTxDone)
			},
			expectedID:    0,
			expectedError: sql.ErrTxDone,
		},
		{
			name: "LastInsertId Error",
			input: &domain.UserProfile{
				UserID:       1,
				NIK:          "123456789",
				FullName:     "John Doe",
				LegalName:    "Johnathan Doe",
				PlaceOfBirth: "City A",
				DateOfBirth:  "1990-01-01",
				Salary:       5000,
				KTP:          "ktp123",
				Selfie:       "selfie123",
			},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT INTO user_profile").
					ExpectExec().
					WithArgs(1, "123456789", "John Doe", "Johnathan Doe", "City A", "1990-01-01", 5000, "ktp123", "selfie123").
					WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))
			},
			expectedID:    0,
			expectedError: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.mockSetup()

			// Call the method
			id, err := repo.AddUserProfile(context.Background(), tt.input)

			// Assert results
			assert.Equal(t, tt.expectedID, id)
			assert.Equal(t, tt.expectedError, err)

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
func TestReadUserProfileByUserID(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Create a new repository instance
	repo := mysql.NewMysqlRepositoryUserProfile(db)

	// Define test cases
	tests := []struct {
		name          string
		userID        int
		mockSetup     func()
		expectedUser  domain.UserProfile
		expectedError error
	}{
		{
			name:   "Success",
			userID: 1,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "nik", "full_name", "legal_name", "place_of_birth", "date_of_birth", "salary", "ktp", "selfie"}).
					AddRow(1, 1, "123456789", "John Doe", "Johnathan Doe", "City A", "1990-01-01", 5000, "ktp123", "selfie123")
				mock.ExpectPrepare("SELECT id, user_id, nik, full_name, legal_name, place_of_birth, date_of_birth, salary, ktp, selfie FROM user_profile WHERE user_id = ?").
					ExpectQuery().
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedUser: domain.UserProfile{
				ID:           1,
				UserID:       1,
				NIK:          "123456789",
				FullName:     "John Doe",
				LegalName:    "Johnathan Doe",
				PlaceOfBirth: "City A",
				DateOfBirth:  "1990-01-01",
				Salary:       5000,
				KTP:          "ktp123",
				Selfie:       "selfie123",
			},
			expectedError: nil,
		},
		{
			name:   "No Rows Found",
			userID: 2,
			mockSetup: func() {
				mock.ExpectPrepare("SELECT id, user_id, nik, full_name, legal_name, place_of_birth, date_of_birth, salary, ktp, selfie FROM user_profile WHERE user_id = ?").
					ExpectQuery().
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows(nil))
			},
			expectedUser:  domain.UserProfile{},
			expectedError: nil,
		},
		{
			name:   "PrepareContext Error",
			userID: 3,
			mockSetup: func() {
				mock.ExpectPrepare("SELECT id, user_id, nik, full_name, legal_name, place_of_birth, date_of_birth, salary, ktp, selfie FROM user_profile WHERE user_id = ?").
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  domain.UserProfile{},
			expectedError: sql.ErrConnDone,
		},
		{
			name:   "QueryRowContext Scan Error",
			userID: 4,
			mockSetup: func() {
				mock.ExpectPrepare("SELECT id, user_id, nik, full_name, legal_name, place_of_birth, date_of_birth, salary, ktp, selfie FROM user_profile WHERE user_id = ?").
					ExpectQuery().
					WithArgs(4).
					WillReturnError(fmt.Errorf("QueryRowContext Scan Error"))
			},
			expectedUser:  domain.UserProfile{},
			expectedError: fmt.Errorf("QueryRowContext Scan Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.mockSetup()

			// Call the method
			user, err := repo.ReadUserProfileByUserID(context.Background(), tt.userID)

			// Assert results
			assert.Equal(t, tt.expectedUser, user)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			}

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
