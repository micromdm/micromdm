package account

import (
	"net/http"

	"github.com/gorilla/csrf"

	"micromdm.io/v2/pkg/frontend"
	"micromdm.io/v2/pkg/log"
)

func (srv server) registerForm(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		logger   = log.FromContext(ctx)
		username = r.FormValue("username")
		email    = r.FormValue("email")
		password = r.FormValue("password")
		data     = frontend.Data{
			csrf.TemplateTag: csrf.TemplateField(r),
			"form": map[string]string{
				"username": username,
				"email":    email,
				"password": password,
			},
		}
	)

	if r.Method == http.MethodGet {
		srv.http.RenderTemplate(ctx, w, "register.tmpl", data)
		return
	}

	usr, err := srv.userdb.CreateUser(ctx, username, email, password)
	if err != nil {
		srv.http.Fail(ctx, w, err, "register.tmpl", "msg", "creating user")
		return
	}

	log.Debug(logger).Log(
		"msg", "account created",
		"username", usr.Username,
		"id", usr.ID,
		"confirmation_hash", *usr.ConfirmationHash,
	)

	http.Redirect(w, r, "/register/done", http.StatusFound)
}

func (srv server) registerComplete(w http.ResponseWriter, r *http.Request) {
	srv.http.RenderTemplate(r.Context(), w, "register-done.tmpl", frontend.Data{})
}
