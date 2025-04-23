package repository_test

import (
	"context"
	"errors"
	"testing"

	"app/domain"
	mysql "app/internal/repository/mysql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddLimitType(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := mysql.NewMysqlRepositoryLimitType(db)

	tests := []struct {
		name          string
		prepareMock   func()
		expectedID    int
		expectedError bool
	}{
		{
			name: "Success",
			prepareMock: func() {
				query := "INSERT limit_type SET tenor = ?"
				mock.ExpectPrepare(query).
					ExpectExec().
					WithArgs(12).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedID:    1,
			expectedError: false,
		},
		{
			name: "Prepare Error",
			prepareMock: func() {
				query := "INSERT limit_type SET tenor = ?"
				mock.ExpectPrepare(query).WillReturnError(errors.New("prepare error"))
			},
			expectedID:    0,
			expectedError: true,
		},
		{
			name: "Exec Error",
			prepareMock: func() {
				query := "INSERT limit_type SET tenor = ?"
				mock.ExpectPrepare(query).
					ExpectExec().
					WithArgs(12).
					WillReturnError(errors.New("exec error"))
			},
			expectedID:    0,
			expectedError: true,
		},
		{
			name: "LastInsertId Error",
			prepareMock: func() {
				query := "INSERT limit_type SET tenor = ?"
				mock.ExpectPrepare(query).
					ExpectExec().
					WithArgs(12).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert id error")))
			},
			expectedID:    0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			limitType := &domain.LimitType{Tenor: 12}

			tt.prepareMock()

			id, err := repo.AddLimitType(ctx, limitType)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedID, id)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestReadLimitTypes(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := mysql.NewMysqlRepositoryLimitType(db)

	tests := []struct {
		name          string
		prepareMock   func()
		expectedLen   int
		expectedError bool
	}{
		{
			name: "Success",
			prepareMock: func() {
				query := "SELECT id, tenor FROM limit_type"
				rows := sqlmock.NewRows([]string{"id", "tenor"}).
					AddRow(1, 12).
					AddRow(2, 24)
				mock.ExpectPrepare(query).
					ExpectQuery().
					WillReturnRows(rows)
			},
			expectedLen:   2,
			expectedError: false,
		},
		{
			name: "Prepare Error",
			prepareMock: func() {
				query := "SELECT id, tenor FROM limit_type"
				mock.ExpectPrepare(query).WillReturnError(errors.New("prepare error"))
			},
			expectedLen:   0,
			expectedError: true,
		},
		{
			name: "Query Error",
			prepareMock: func() {
				query := "SELECT id, tenor FROM limit_type"
				mock.ExpectPrepare(query).
					ExpectQuery().
					WillReturnError(errors.New("query error"))
			},
			expectedLen:   0,
			expectedError: true,
		},
		{
			name: "Scan Error",
			prepareMock: func() {
				query := "SELECT id, tenor FROM limit_type"
				rows := sqlmock.NewRows([]string{"id", "tenor"}).
					AddRow(1, "invalid_tenor")
				mock.ExpectPrepare(query).
					ExpectQuery().
					WillReturnRows(rows)
			},
			expectedLen:   0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tt.prepareMock()

			limitTypes, err := repo.ReadLimitTypes(ctx)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, limitTypes)
			} else {
				assert.NoError(t, err)
				assert.Len(t, limitTypes, tt.expectedLen)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
