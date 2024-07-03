package auth

import (
	"context"
	"errors"
	"testing"

	"todo/internal/models"
	"todo/pkg/log"
	pb "todo/pkg/proto/auth"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthGRPC_Register(t *testing.T) {
	tests := []struct {
		name          string
		login         string
		password      string
		mockSetup     func(*MockAuthService)
		expectedError error
	}{
		{
			name:     "success",
			login:    "testuser",
			password: "password",
			mockSetup: func(m *MockAuthService) {
				m.On("Register", "testuser", "password").Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "user already exists",
			login:    "existinguser",
			password: "password",
			mockSetup: func(m *MockAuthService) {
				m.On("Register", "existinguser", "password").Return(models.ErrUserExists)
			},
			expectedError: status.Error(codes.AlreadyExists, models.ErrUserExists.Error()),
		},
		{
			name:     "internal error",
			login:    "testuser",
			password: "password",
			mockSetup: func(m *MockAuthService) {
				m.On("Register", "testuser", "password").Return(errors.New("some internal error"))
			},
			expectedError: status.Error(codes.Internal, "internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := new(MockAuthService)
			mockLogger := new(log.MockLogger)
			authServer := AuthServer{service: mockAuthService, logger: mockLogger}

			tt.mockSetup(mockAuthService)

			_, err := authServer.Register(context.Background(), &pb.RegisterRequest{Login: tt.login, Password: tt.password})

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockAuthService.AssertExpectations(t)
		})
	}
}
