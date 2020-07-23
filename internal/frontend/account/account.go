// Package account contains web pages for user registration and account management features.
package account

import (
	"context"
	"net/http"

	"micromdm.io/v2/internal/data/user"
	"micromdm.io/v2/pkg/frontend"
)

type UserStore interface {
	CreateUser(ctx context.Context, username, email, password string) (*user.User, error)
	ConfirmUser(ctx context.Context, token string) error
}

type server struct {
	http   frontend.Framework
	userdb UserStore
}

type Config struct {
	HTTP      frontend.Framework
	UserStore UserStore
}

func HTTP(config Config) {
	srv := &server{
		http:   config.HTTP,
		userdb: config.UserStore,
	}

	srv.http.HandleFunc("/register", srv.registerForm, http.MethodGet, http.MethodPost)
	srv.http.HandleFunc("/register/done", srv.registerComplete)
	srv.http.HandleFunc("/registered/confirm/{token}", srv.registerConfirm)
}
