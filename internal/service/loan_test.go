package service_test

import (
	"app/domain"
	"app/internal/service"
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepositoryMySQLLoan is a mock implementation of the RepositoryMySQLLoan interface
type MockRepositoryMySQLLoan struct {
	mock.Mock
}

func (m *MockRepositoryMySQLLoan) BeginTx(ctx context.Context) (*sql.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(*sql.Tx), args.Error(1)
}

func (m *MockRepositoryMySQLLoan) AddTXLoan(ctx context.Context, tx *sql.Tx, lo *domain.Loan) (int, error) {
	args := m.Called(ctx, tx, lo)
	return args.Int(0), args.Error(1)
}

// MockRepositoryMySQLLimit is a mock implementation of the RepositoryMySQLLimit interface
type MockRepositoryMySQLLimit struct {
	mock.Mock
}

func (m *MockRepositoryMySQLLimit) ReadTXLimitByID(ctx context.Context, tx *sql.Tx, id int) (domain.Limit, error) {
	args := m.Called(ctx, tx, id)
	limit, _ := args.Get(0).(*domain.Limit)
	return *limit, args.Error(1)
}

func (m *MockRepositoryMySQLLimit) UpdateTXLimit(ctx context.Context, tx *sql.Tx, limit *domain.Limit) error {
	args := m.Called(ctx, tx, limit)
	return args.Error(0)
}

func TestServiceAddLoan(t *testing.T) {
	mockLoanRepo := new(MockRepositoryMySQLLoan)
	mockLimitRepo := new(MockRepositoryMySQLLimit)
	service := service.NewServiceLoan(mockLoanRepo, mockLimitRepo)

	tests := []struct {
		name          string
		input         *domain.Loan
		mockBehavior  func()
		expectedID    int
		expectedError bool
	}{
		{
			name: "Success",
			input: &domain.Loan{
				LimitID: 1,
				Amount:  500,
			},
			mockBehavior: func() {
				db, mck, err := sqlmock.New()
				assert.NoError(t, err)
				defer db.Close()

				mck.ExpectBegin()

				tx, _ := db.Begin()

				mockLoanRepo.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				mockLimitRepo.On("ReadTXLimitByID", mock.Anything, tx, 1).Return(&domain.Limit{
					ID:         1,
					Amount:     1000,
					UsedAmount: 400,
				}, nil).Once()
				mockLoanRepo.On("AddTXLoan", mock.Anything, tx, mock.Anything).Return(1, nil).Once()
				mockLimitRepo.On("UpdateTXLimit", mock.Anything, tx, &domain.Limit{
					ID:         1,
					UsedAmount: 900,
				}).Return(nil).Once()
			},
			expectedID:    1,
			expectedError: false,
		},
		{
			name: "Limit Exceeded",
			input: &domain.Loan{
				LimitID: 1,
				Amount:  700,
			},
			mockBehavior: func() {
				db, mck, err := sqlmock.New()
				assert.NoError(t, err)
				defer db.Close()

				mck.ExpectBegin()

				tx, _ := db.Begin()

				mockLoanRepo.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				mockLimitRepo.On("ReadTXLimitByID", mock.Anything, tx, 1).Return(&domain.Limit{
					ID:         1,
					Amount:     1000,
					UsedAmount: 400,
				}, nil).Once()
			},
			expectedID:    0,
			expectedError: true,
		},
		{
			name: "Repository Error",
			input: &domain.Loan{
				LimitID: 1,
				Amount:  500,
			},
			mockBehavior: func() {
				db, mck, err := sqlmock.New()
				assert.NoError(t, err)
				defer db.Close()

				mck.ExpectBegin()

				tx, _ := db.Begin()

				mockLoanRepo.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				mockLimitRepo.On("ReadTXLimitByID", mock.Anything, tx, 1).Return(&domain.Limit{
					ID:         1,
					Amount:     1000,
					UsedAmount: 400,
				}, nil).Once()
				mockLoanRepo.On("AddTXLoan", mock.Anything, tx, mock.Anything).Return(0, errors.New("repository error")).Once()
			},
			expectedID:    0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			id, err := service.AddLoan(context.Background(), tt.input)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedID, id)
			mockLoanRepo.AssertExpectations(t)
			mockLimitRepo.AssertExpectations(t)
		})
	}
}
