package todoapp

import (
	"context"
	"net/http"
	"todo/internal/jwt"
	"todo/internal/models"
)

type LoginKey string

func (a *App) ProtectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			a.handleError(w, models.ErrUnauthorized)
			return
		}

		login, err := jwt.ParseJWT(token, a.config.JwtSecret)
		if err != nil {
			a.handleError(w, models.ErrUnauthorized)
			return
		}

		const loginKey LoginKey = "login"
		ctx := context.WithValue(r.Context(), loginKey, login)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
