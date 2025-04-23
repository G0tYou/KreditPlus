package service_test

import (
	"app/domain"
	"app/internal/service"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepositoryMySQLUser is a mock implementation of the RepositoryMySQLUser interface
type MockRepositoryMySQLUser struct {
	mock.Mock
}

func (m *MockRepositoryMySQLUser) ExistByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepositoryMySQLUser) AddUser(ctx context.Context, u *domain.User) (int, error) {
	args := m.Called(ctx, u)
	return args.Int(0), args.Error(1)
}

func (m *MockRepositoryMySQLUser) ReadUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return &domain.User{}, args.Error(1)
	}
	user, _ := args.Get(0).(domain.User)
	return &user, args.Error(1)
}

func TestServiceAddUser(t *testing.T) {
	mockRepo := new(MockRepositoryMySQLUser)
	service := service.NewServiceUser(mockRepo)

	tests := []struct {
		name          string
		input         *domain.User
		mockBehavior  func()
		expectedID    int
		expectedError bool
	}{
		{
			name: "Success",
			input: &domain.User{
				Username: "testuser",
				Password: "password123",
			},
			mockBehavior: func() {
				mockRepo.On("ExistByUsername", mock.Anything, "testuser").Return(false, nil).Once()
				mockRepo.On("AddUser", mock.Anything, mock.Anything).Return(1, nil).Once()
			},
			expectedID:    1,
			expectedError: false,
		},
		{
			name: "Username Exists",
			input: &domain.User{
				Username: "testuser",
				Password: "password123",
			},
			mockBehavior: func() {
				mockRepo.On("ExistByUsername", mock.Anything, "testuser").Return(true, nil).Once()
			},
			expectedID:    0,
			expectedError: true,
		},
		{
			name: "Repository Error",
			input: &domain.User{
				Username: "testuser",
				Password: "password123",
			},
			mockBehavior: func() {
				mockRepo.On("ExistByUsername", mock.Anything, "testuser").Return(false, nil).Once()
				mockRepo.On("AddUser", mock.Anything, mock.Anything).Return(0, errors.New("repository error")).Once()
			},
			expectedID:    0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			id, err := service.AddUser(context.Background(), tt.input)

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

func TestServiceLogin(t *testing.T) {
	mockRepo := new(MockRepositoryMySQLUser)
	service := service.NewServiceUser(mockRepo)

	tests := []struct {
		name          string
		input         *domain.User
		mockBehavior  func()
		expectedError bool
	}{
		{
			name: "Success",
			input: &domain.User{
				Username: "dachiw",
				Password: "dachiw",
			},
			mockBehavior: func() {
				mockRepo.On("ExistByUsername", mock.Anything, "dachiw").Return(true, nil).Once()
				mockRepo.On("ReadUserByUsername", mock.Anything, "dachiw").Return(domain.User{
					ID:       1,
					Username: "dachiw",
					Password: "$2a$10$Lff.n7y8piIAIrt1V.qC9Oj8Ico4nU7T3KPIiCvvtBJaVP20Qorvy",
				}, nil).Once()
			},
			expectedError: false,
		},
		{
			name: "User Not Found",
			input: &domain.User{
				Username: "testuser",
				Password: "password123",
			},
			mockBehavior: func() {
				mockRepo.On("ExistByUsername", mock.Anything, "testuser").Return(false, nil).Once()
			},
			expectedError: true,
		},
		{
			name: "Invalid Password",
			input: &domain.User{
				Username: "testuser",
				Password: "wrongpassword",
			},
			mockBehavior: func() {
				mockRepo.On("ExistByUsername", mock.Anything, "testuser").Return(true, nil).Once()
				mockRepo.On("ReadUserByUsername", mock.Anything, "testuser").Return(domain.User{
					ID:       1,
					Username: "testuser",
					Password: "hashedpassword",
				}, nil).Once()
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			_, err := service.Login(context.Background(), tt.input)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
