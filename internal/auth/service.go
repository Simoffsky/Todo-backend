package auth

import (
	"time"
	"todo/internal/models"
	repository "todo/internal/repository/auth"
	"todo/pkg/log"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService interface {
	Register(login, password string) error
	Login(login, password string) (string, error)
	UserExists(login string) (bool, error)
}

type AuthServiceDefault struct {
	logger    log.Logger
	repo      repository.UserRepository
	jwtSecret string
}

func NewDefaultAuthService(logger log.Logger, userRepo repository.UserRepository, jwtSecret string) *AuthServiceDefault {
	return &AuthServiceDefault{
		logger:    logger,
		repo:      userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthServiceDefault) Register(login, password string) error {
	exists, err := s.repo.Exists(login)
	if err != nil {
		return err
	}
	if exists {
		return models.ErrUserExists
	}

	user := &models.User{
		Login:    login,
		Password: password,
	}
	s.logger.Debug("Registered user: " + login)
	return s.repo.CreateUser(user)
}

func (s *AuthServiceDefault) Login(login, password string) (string, error) {
	user, err := s.repo.GetUser(login)
	if err != nil {
		return "", err
	}

	if user.Password != password {
		return "", models.ErrInvalidPassword
	}

	payload := jwt.MapClaims{
		"sub": user.Login,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	s.logger.Debug("User logged in: " + login)
	return signedToken, nil
}

func (s *AuthServiceDefault) UserExists(login string) (bool, error) {
	return s.repo.Exists(login)
}
