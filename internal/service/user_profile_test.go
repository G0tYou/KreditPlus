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

// MockRepositoryMySQLUserProfile is a mock implementation of the RepositoryMySQLUserProfile interface
type MockRepositoryMySQLUserProfile struct {
	mock.Mock
}

func (m *MockRepositoryMySQLUserProfile) AddUserProfile(ctx context.Context, up *domain.UserProfile) (int, error) {
	args := m.Called(ctx, up)
	return args.Int(0), args.Error(1)
}

func (m *MockRepositoryMySQLUserProfile) ReadUserProfileByUserID(ctx context.Context, userID int) (domain.UserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return domain.UserProfile{}, args.Error(1)
	}
	return args.Get(0).(domain.UserProfile), args.Error(1)
}

func TestServiceAddUserProfile(t *testing.T) {
	mockRepo := new(MockRepositoryMySQLUserProfile)

	service := service.NewServiceUserProfile(mockRepo)

	tests := []struct {
		name          string
		input         *domain.UserProfile
		mockBehavior  func()
		expectedID    int
		expectedError bool
	}{
		{
			name: "Success",
			input: &domain.UserProfile{
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
			mockBehavior: func() {
				mockRepo.On("AddUserProfile", mock.Anything, mock.Anything).Return(1, nil).Once()
			},
			expectedID:    1,
			expectedError: false,
		},
		{
			name: "Repository Error",
			input: &domain.UserProfile{
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
			mockBehavior: func() {
				mockRepo.On("AddUserProfile", mock.Anything, mock.Anything).Return(0, errors.New("repository error")).Once()
			},
			expectedID:    0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			id, err := service.AddUserProfile(context.Background(), tt.input)

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

func TestServiceGetUserProfileByUserID(t *testing.T) {
	mockRepo := new(MockRepositoryMySQLUserProfile)

	service := service.NewServiceUserProfile(mockRepo)

	tests := []struct {
		name          string
		userID        int
		mockBehavior  func()
		expectedData  domain.UserProfile
		expectedError bool
	}{
		{
			name:   "Success",
			userID: 1,
			mockBehavior: func() {
				mockRepo.On("ReadUserProfileByUserID", mock.Anything, 1).
					Return(domain.UserProfile{
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
					}, nil).Once()
			},
			expectedData: domain.UserProfile{
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
			expectedError: false,
		},
		{
			name:   "Repository Error",
			userID: 1,
			mockBehavior: func() {
				mockRepo.On("ReadUserProfileByUserID", mock.Anything, 1).
					Return(domain.UserProfile{}, errors.New("repository error")).Once()
			},
			expectedData:  domain.UserProfile{},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			data, err := service.GetUserProfileByUserID(context.Background(), tt.userID)

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
