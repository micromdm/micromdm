// Package account contains web pages for user registration and account management features.
package account

import (
	"context"
	"net/http"

	"github.com/gorilla/securecookie"
	"micromdm.io/v2/internal/data/session"
	"micromdm.io/v2/internal/data/user"
	"micromdm.io/v2/pkg/frontend"
)

type UserStore interface {
	CreateUser(ctx context.Context, username, email, password string) (*user.User, error)
	ConfirmUser(ctx context.Context, token string) error
	FindUserByEmail(ctx context.Context, email string) (*user.User, error)
}

type SessionStore interface {
	CreateSession(ctx context.Context) (*session.Session, error)
	DestroySession(ctx context.Context) error
}

type CookieAuthFramework interface {
	frontend.Framework
	AuthCookieName() string
}

type server struct {
	http      CookieAuthFramework
	userdb    UserStore
	sessiondb SessionStore
	cookie    *securecookie.SecureCookie
}

type Config struct {
	HTTP         CookieAuthFramework
	UserStore    UserStore
	SessionStore SessionStore
	Cookie       *securecookie.SecureCookie
}

func HTTP(config Config) {
	srv := &server{
		http:      config.HTTP,
		userdb:    config.UserStore,
		sessiondb: config.SessionStore,
		cookie:    config.Cookie,
	}

	srv.http.HandleFunc("/register", srv.registerForm, http.MethodGet, http.MethodPost)
	srv.http.HandleFunc("/register/done", srv.registerComplete)
	srv.http.HandleFunc("/registered/confirm/{token}", srv.registerConfirm)

	srv.http.HandleFunc("/login", srv.loginForm, http.MethodGet, http.MethodPost)
	srv.http.HandleFunc("/logout", srv.logout)
}
