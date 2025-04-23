package repository_test

import (
	"context"
	"database/sql"
	"testing"

	d "app/domain"
	mysql "app/internal/repository/mysql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddTXLoan(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := mysql.NewMysqlRepositoryLoan(db)

	mock.ExpectBegin()

	tx, _ := db.BeginTx(context.Background(), &sql.TxOptions{})

	tests := []struct {
		name          string
		Loan          d.Loan
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Success",
			Loan: d.Loan{
				LimitID: 1,
				Amount:  1000,
			},
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO loan").
					WithArgs(1, 1000).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: nil,
		},
		{
			name: "Error Insert",
			Loan: d.Loan{
				LimitID: 2,
				Amount:  2000,
			},
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO loan").
					WithArgs(2, 2000).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			_, err := repo.AddTXLoan(context.Background(), tx, &tt.Loan)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			}

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
