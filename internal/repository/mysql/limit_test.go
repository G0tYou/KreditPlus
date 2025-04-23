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

func TestReadTXLimitByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := mysql.NewMysqlRepositoryLimit(db)

	mock.ExpectBegin()

	tx, _ := db.BeginTx(context.Background(), &sql.TxOptions{})

	tests := []struct {
		name          string
		limitID       int
		mockSetup     func()
		expectedLimit d.Limit
		expectedError error
	}{
		{
			name:    "Success",
			limitID: 1,
			mockSetup: func() {
				query := "SELECT id, limit_type_id, user_id, amount, used_amount FROM `limit` WHERE id = \\? FOR UPDATE"
				rows := sqlmock.NewRows([]string{"id", "limit_type_id", "user_id", "amount", "used_amount"}).
					AddRow(1, 2, 3, 1000, 500)

				mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
			},
			expectedLimit: d.Limit{
				ID:          1,
				LimitTypeID: 2,
				UserID:      3,
				Amount:      1000,
				UsedAmount:  500,
			},
			expectedError: nil,
		},
		{
			name:    "Not Rows Found",
			limitID: 1,
			mockSetup: func() {
				query := "SELECT id, limit_type_id, user_id, amount, used_amount FROM `limit` WHERE id = \\? FOR UPDATE"
				mock.ExpectQuery(query).WithArgs(1).WillReturnError(sql.ErrNoRows)
			},
			expectedLimit: d.Limit{},
			expectedError: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			limit, err := repo.ReadTXLimitByID(context.Background(), tx, tt.limitID)

			assert.Equal(t, tt.expectedLimit, limit)
			assert.Equal(t, tt.expectedError, err)

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateTXLimit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := mysql.NewMysqlRepositoryLimit(db)

	mock.ExpectBegin()

	tx, _ := db.BeginTx(context.Background(), &sql.TxOptions{})

	tests := []struct {
		name          string
		limit         *d.Limit
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Success",
			limit: &d.Limit{
				ID:         1,
				UsedAmount: 600,
			},
			mockSetup: func() {
				query := "UPDATE `limit` SET used_amount = \\? WHERE id = \\?"
				mock.ExpectExec(query).WithArgs(600, 1).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name: "Error",
			limit: &d.Limit{
				ID:         1,
				UsedAmount: 600,
			},
			mockSetup: func() {
				query := "UPDATE `limit` SET used_amount = \\? WHERE id = \\?"
				mock.ExpectExec(query).WithArgs(600, 1).WillReturnError(sql.ErrConnDone)
			},
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.UpdateTXLimit(context.Background(), tx, tt.limit)

			assert.Equal(t, tt.expectedError, err)

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
