package auth

import (
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(login, password string) error {
	args := m.Called(login, password)
	return args.Error(0)
}

func (m *MockAuthService) Login(login, password string) (string, error) {
	args := m.Called(login, password)
	return args.String(0), args.Error(1)
}
