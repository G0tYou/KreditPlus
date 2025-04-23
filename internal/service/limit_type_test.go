package service_test

import (
	"app/domain"
	"context"
	"errors"
	"testing"

	"app/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepositoryMySQLLimitType is a mock implementation of the RepositoryMySQLLimitType interface
type MockRepositoryMySQLLimitType struct {
	mock.Mock
}

func (m *MockRepositoryMySQLLimitType) AddLimitType(ctx context.Context, lt *domain.LimitType) (int, error) {
	args := m.Called(ctx, lt)
	return args.Int(0), args.Error(1)
}

func (m *MockRepositoryMySQLLimitType) ReadLimitTypes(ctx context.Context) ([]domain.LimitType, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.LimitType), args.Error(1)
}

func TestServiceAddLimitType(t *testing.T) {
	mockRepo := new(MockRepositoryMySQLLimitType)

	service := service.NewServiceLimitType(mockRepo)

	tests := []struct {
		name          string
		input         *domain.LimitType
		mockBehavior  func()
		expectedID    int
		expectedError bool
	}{
		{
			name: "Success",
			input: &domain.LimitType{
				Tenor: 12,
			},
			mockBehavior: func() {
				mockRepo.On("AddLimitType", mock.Anything, &domain.LimitType{Tenor: 12}).
					Return(1, nil).Once()
			},
			expectedID:    1,
			expectedError: false,
		},
		{
			name: "Repository Error",
			input: &domain.LimitType{
				Tenor: 12,
			},
			mockBehavior: func() {
				mockRepo.On("AddLimitType", mock.Anything, &domain.LimitType{Tenor: 12}).
					Return(0, errors.New("repository error")).Once()
			},
			expectedID:    0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			id, err := service.AddLimitType(context.Background(), tt.input)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedID, id)
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestServiceGetLimitTypes(t *testing.T) {
	mockRepo := new(MockRepositoryMySQLLimitType)

	service := service.NewServiceLimitType(mockRepo)

	tests := []struct {
		name          string
		mockBehavior  func()
		expectedData  []domain.LimitType
		expectedError bool
	}{
		{
			name: "Success",
			mockBehavior: func() {
				mockRepo.On("ReadLimitTypes", mock.Anything).
					Return([]domain.LimitType{
						{ID: 1, Tenor: 12},
						{ID: 2, Tenor: 24},
					}, nil).Once()
			},
			expectedData: []domain.LimitType{
				{ID: 1, Tenor: 12},
				{ID: 2, Tenor: 24},
			},
			expectedError: false,
		},
		{
			name: "Repository Error",
			mockBehavior: func() {
				mockRepo.On("ReadLimitTypes", mock.Anything).
					Return(nil, errors.New("repository error")).Once()
			},
			expectedData:  nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			data, err := service.GetLimitTypes(context.Background())

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedData, data)
			mockRepo.AssertExpectations(t)
		})
	}
}
