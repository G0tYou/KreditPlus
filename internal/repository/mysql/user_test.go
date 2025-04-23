package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"app/domain"
	mysql "app/internal/repository/mysql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := mysql.NewMysqlRepositoryUser(db)

	tests := []struct {
		name          string
		inputUser     *domain.User
		mockSetup     func()
		expectedID    int
		expectedError error
	}{
		{
			name: "Success",
			inputUser: &domain.User{
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT user SET ").
					ExpectExec().
					WithArgs("testuser", "password123").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedID:    1,
			expectedError: nil,
		},
		{
			name: "Prepare Error",
			inputUser: &domain.User{
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT user SET ").
					WillReturnError(sql.ErrConnDone)
			},
			expectedID:    0,
			expectedError: sql.ErrConnDone,
		},
		{
			name: "Exec Error",
			inputUser: &domain.User{
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT user SET ").
					ExpectExec().
					WithArgs("testuser", "password123").
					WillReturnError(errors.New("exec error"))
			},
			expectedID:    0,
			expectedError: errors.New("exec error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			id, err := repo.AddUser(context.Background(), tt.inputUser)

			assert.Equal(t, tt.expectedID, id)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			}

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestReadUserByUsername(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		mockSetup     func(mock sqlmock.Sqlmock)
		wantUser      *domain.User
		expectedError error
	}{
		{
			name:     "Success",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "password"}).
					AddRow(1, "testuser", "password123")
				mock.ExpectPrepare("SELECT id, username, password FROM user WHERE username = ?").
					ExpectQuery().
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			wantUser: &domain.User{
				ID:       1,
				Username: "testuser",
				Password: "password123",
			},
			expectedError: nil,
		},
		{
			name:     "No Rows Found",
			username: "nonexistent",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("SELECT id, username, password FROM user WHERE username = ?").
					ExpectQuery().
					WithArgs("nonexistent").
					WillReturnError(sql.ErrNoRows)
			},
			wantUser:      nil,
			expectedError: domain.ErrNotFound,
		},
		{
			name:     "Query Error",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("SELECT id, username, password FROM user WHERE username = ?").
					ExpectQuery().
					WithArgs("testuser").
					WillReturnError(errors.New("query error"))
			},
			wantUser:      nil,
			expectedError: errors.New("query error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock db: %v", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			repo := mysql.NewMysqlRepositoryUser(db)
			user, err := repo.ReadUserByUsername(context.Background(), tt.username)

			if user != nil && tt.wantUser != nil {
				if user.ID != tt.wantUser.ID || user.Username != tt.wantUser.Username || user.Password != tt.wantUser.Password {
					t.Errorf("got user %+v, want %+v", user, tt.wantUser)
				}
			} else if user != tt.wantUser {
				t.Errorf("got user %+v, want %+v", user, tt.wantUser)
			}

			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("got error %v, want %v", err, tt.expectedError)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet mock expectations: %v", err)
			}
		})
	}
}
func TestExistByUsername(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		mockSetup     func(mock sqlmock.Sqlmock)
		wantExist     bool
		expectedError error
	}{
		{
			name:     "User Exists",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"EXISTS(SELECT 1 FROM user WHERE username = ?)"}).
					AddRow(true)
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM user WHERE username = \\?\\)").
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			wantExist:     true,
			expectedError: nil,
		},
		{
			name:     "User Does Not Exist",
			username: "nonexistent",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"EXISTS(SELECT 1 FROM user WHERE username = ?)"}).
					AddRow(false)
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM user WHERE username = \\?\\)").
					WithArgs("nonexistent").
					WillReturnRows(rows)
			},
			wantExist:     false,
			expectedError: nil,
		},
		{
			name:     "Query Error",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM user WHERE username = \\?\\)").
					WithArgs("testuser").
					WillReturnError(errors.New("query error"))
			},
			wantExist:     false,
			expectedError: errors.New("query error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock db: %v", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			repo := mysql.NewMysqlRepositoryUser(db)
			exists, err := repo.ExistByUsername(context.Background(), tt.username)

			if exists != tt.wantExist {
				t.Errorf("got exists %v, want %v", exists, tt.wantExist)
			}
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("got error %v, want %v", err, tt.expectedError)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet mock expectations: %v", err)
			}
		})
	}
}
