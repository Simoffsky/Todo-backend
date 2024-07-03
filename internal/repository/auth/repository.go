package auth

import "todo/internal/models"

type UserRepository interface {
	CreateUser(*models.User) error
	Exists(login string) (bool, error)
	GetUser(login string) (*models.User, error)
}
