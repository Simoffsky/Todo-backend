package auth

import (
	"errors"
	"testing"
	"todo/internal/models"

	repository "todo/internal/repository/auth"
	"todo/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name          string
		login         string
		password      string
		mockSetup     func(*repository.MockUserRepository)
		expectedError error
	}{
		{
			name:     "success",
			login:    "user@example.com",
			password: "password",
			mockSetup: func(m *repository.MockUserRepository) {
				m.On("Exists", "user@example.com").Return(false, nil)
				m.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "user already exists",
			login:    "user@example.com",
			password: "password",
			mockSetup: func(m *repository.MockUserRepository) {
				m.On("Exists", "user@example.com").Return(true, nil)
			},
			expectedError: models.ErrUserExists,
		},
		{
			name:     "repository error on exists check",
			login:    "user@example.com",
			password: "password",
			mockSetup: func(m *repository.MockUserRepository) {
				m.On("Exists", "user@example.com").Return(false, errors.New("unexpected error"))
			},
			expectedError: errors.New("unexpected error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockUserRepository)
			logger := &log.MockLogger{}
			service := NewDefaultAuthService(logger, mockRepo, "secret")

			tt.mockSetup(mockRepo)

			err := service.Register(tt.login, tt.password)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name          string
		login         string
		password      string
		mockSetup     func(*repository.MockUserRepository)
		expectedError error
		expectedToken bool
	}{
		{
			name:     "success",
			login:    "user@example.com",
			password: "password",
			mockSetup: func(m *repository.MockUserRepository) {
				m.On("GetUser", "user@example.com").Return(&models.User{Login: "user@example.com", Password: "password"}, nil)
			},
			expectedError: nil,
			expectedToken: true,
		},
		{
			name:     "invalid password",
			login:    "user@example.com",
			password: "wrongpassword",
			mockSetup: func(m *repository.MockUserRepository) {
				m.On("GetUser", "user@example.com").Return(&models.User{Login: "user@example.com", Password: "password"}, nil)
			},
			expectedError: models.ErrInvalidPassword,
			expectedToken: false,
		},
		{
			name:     "user not found",
			login:    "user@example.com",
			password: "password",
			mockSetup: func(m *repository.MockUserRepository) {
				m.On("GetUser", "user@example.com").Return(nil, models.ErrUserNotFound)
			},
			expectedError: models.ErrUserNotFound,
			expectedToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockUserRepository)
			logger := &log.MockLogger{}
			service := NewDefaultAuthService(logger, mockRepo, "secret")

			tt.mockSetup(mockRepo)

			token, err := service.Login(tt.login, tt.password)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				if tt.expectedToken {
					assert.NotEmpty(t, token)
				} else {
					assert.Empty(t, token)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_UserExists(t *testing.T) {
	tests := []struct {
		name           string
		login          string
		mockSetup      func(*repository.MockUserRepository)
		expectedError  error
		expectedExists bool
	}{
		{
			name:  "user exists",
			login: "existing",
			mockSetup: func(m *repository.MockUserRepository) {
				m.On("Exists", "existing").Return(true, nil)
			},
			expectedError:  nil,
			expectedExists: true,
		},
		{
			name:  "user does not exist",
			login: "nonexisting",
			mockSetup: func(m *repository.MockUserRepository) {
				m.On("Exists", "nonexisting").Return(false, nil)
			},
			expectedError:  nil,
			expectedExists: false,
		},
		{
			name:  "repository error",
			login: "aboba",
			mockSetup: func(m *repository.MockUserRepository) {
				m.On("Exists", "aboba").Return(false, errors.New("unexpected error"))
			},
			expectedError:  errors.New("unexpected error"),
			expectedExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockUserRepository)
			logger := &log.MockLogger{}
			service := NewDefaultAuthService(logger, mockRepo, "secret")

			tt.mockSetup(mockRepo)

			exists, err := service.UserExists(tt.login)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedExists, exists)
			}

			mockRepo.AssertExpectations(t)
		})
	}

}
