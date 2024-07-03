package auth

import (
	"todo/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Exists(login string) (bool, error) {
	args := m.Called(login)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) GetUser(login string) (*models.User, error) {
	args := m.Called(login)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*models.User), args.Error(1)
}
