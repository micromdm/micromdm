package account

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"micromdm.io/v2/pkg/frontend"
	"micromdm.io/v2/pkg/log"
	"micromdm.io/v2/pkg/viewer"
)

func (srv server) loginForm(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		logger   = log.FromContext(ctx)
		email    = r.FormValue("email")
		password = r.FormValue("password")
		data     = frontend.Data{
			csrf.TemplateTag: csrf.TemplateField(r),
			"form": map[string]string{
				"email":    email,
				"password": password,
			},
		}
	)

	if r.Method == http.MethodGet {
		srv.http.RenderTemplate(ctx, w, "login.tmpl", data)
		return
	}

	usr, err := srv.userdb.FindUserByEmail(ctx, email)
	if err != nil {
		srv.http.Fail(ctx, w, err, "msg", "find user for auth")
		return
	}

	if err := usr.ValidatePassword(password); err != nil {
		srv.http.Fail(ctx, w, err, "msg", "auth user")
		return
	}

	log.Info(logger).Log("msg", "got user", "user_id", usr.ID)
	if err := srv.createSession(ctx, w, usr.ID); err != nil {
		srv.http.Fail(ctx, w, err, "msg", "create session")
		return
	}
	log.Info(logger).Log("msg", "logged in", "user_id", usr.ID)
}

func (srv server) createSession(ctx context.Context, w http.ResponseWriter, userID string) error {
	ctx = viewer.NewContext(ctx, viewer.Viewer{UserID: userID})
	sess, err := srv.sessiondb.CreateSession(ctx)
	if err != nil {
		return fmt.Errorf("create session for %q: %w", userID, err)
	}

	token, err := srv.cookie.Encode(srv.http.AuthCookieName(), map[string]string{"id": sess.ID})
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     srv.http.AuthCookieName(),
		Path:     "/",
		Value:    token,
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().UTC().Add(30 * time.Minute),
	})

	return nil
}
