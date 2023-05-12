package controllers

import (
	"net/http"

	"github.com/alorents/lenslocked/context"
	"github.com/alorents/lenslocked/models"
)

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := readCooke(r, CookeSession)
		if err != nil || tokenCookie.Value == "" {
			next.ServeHTTP(w, r)
			return
		}

		user, err := umw.SessionService.User(tokenCookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithUser(r.Context(), user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
