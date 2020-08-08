// Package frontend provides a lightweight framework for building the MicroMDM HTML UI.
package frontend

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"

	"micromdm.io/v2/internal/data/session"
	"micromdm.io/v2/pkg/log"
	"micromdm.io/v2/pkg/viewer"
)

// Data keys. Private, set via helper methods.
const (
	dHTTPCode  = "http-code"
	dLogErr    = "log-error"
	dLogKV     = "log-keyvals"
	dFormErrs  = "errors"
	dFormAlert = "alert"
)

// Data provides request parameters when calling RenderTemplate.
type Data map[string]interface{}

// WithLog adds keyvals to log when rendering a template.
func (d Data) WithLog(err error, keyvals ...interface{}) Data {
	d[dLogErr] = err
	d[dLogKV] = keyvals
	return d
}

// WithCode sets an HTTP status code. The default value when not set is 200 OK.
func (d Data) WithCode(code int) Data {
	d[dHTTPCode] = code
	return d
}

// FormErrors adds an "errors" key with a mapping of form fields names to error messages.
// FormErrors sets the HTTP status code to 400 StatusBadRequest.
func (d Data) FormErrors(errs map[string]string) Data {
	d[dFormErrs] = errs
	return d.WithCode(http.StatusBadRequest)
}

// Framework specifies methods frontend sub-packages depend on.
// Framework is mainly exported to give sub-packages a common interface to depend on.
// Server is the only used Framework implementation.
type Framework interface {
	Fail(ctx context.Context, w http.ResponseWriter, err error, keyvals ...interface{})
	RenderTemplate(ctx context.Context, w http.ResponseWriter, name string, data Data)
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request), methods ...string)
}

type SessionStore interface {
	FindSession(ctx context.Context, id string) (*session.Session, error)
}

// Server implements Framework.
type Server struct {
	r *mux.Router

	mu        sync.Mutex
	templates map[string]*template.Template

	siteName string

	csrfKey        []byte
	csrfCookieName string
	csrfFieldName  string

	authCookieName string
	sessiondb      SessionStore
	cookie         *securecookie.SecureCookie
}

// Config parameters to create a new Server.
type Config struct {
	Logger   log.Logger
	SiteName string

	CSRFKey        []byte
	CSRFCookieName string
	CSRFFieldName  string

	AuthCookieName string
	SessionStore   SessionStore
	Cookie         *securecookie.SecureCookie
}

// New creates a Server.
func New(config Config) (*Server, error) {
	srv := &Server{
		r:         mux.NewRouter(),
		templates: make(map[string]*template.Template),
		siteName:  config.SiteName,

		csrfKey:        config.CSRFKey,
		csrfFieldName:  config.CSRFFieldName,
		csrfCookieName: config.CSRFCookieName,
		authCookieName: config.AuthCookieName,
		sessiondb:      config.SessionStore,
		cookie:         config.Cookie,
	}

	srv.r.Use(
		log.HTTP(config.Logger), // HTTP logging middleware.
		srv.authMW,              // check auth, add viewer
		srv.csrf,                // CSRF protection.
		srv.recoverPanic,        // convert any panic into 500 errors.
	)

	// have to set middleware for NotFoundHandler separate from matched routes.
	srv.r.NotFoundHandler = log.HTTP(config.Logger)(http.HandlerFunc(srv.notFound))

	srv.r.HandleFunc("/", srv.indexPage).Methods(http.MethodGet)

	// Serve all static content.
	// This is another place that will need to be improved to serve from a CDN or object store instead.
	srv.r.PathPrefix("/assets/").Methods(http.MethodGet).
		Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-cache")
			http.StripPrefix("/assets/", http.FileServer(http.Dir("ui/static"))).ServeHTTP(w, r)
		}))

	if err := srv.loadTemplates(); err != nil {
		return nil, fmt.Errorf("loading ui templates: %s", err)
	}

	return srv, nil
}

// Handler returns the mux router used by the Server.
func (srv *Server) Handler() http.Handler { return srv.r }

// AuthCookieName exposes the cookie name used for authentication.
func (srv *Server) AuthCookieName() string { return srv.authCookieName }

// HandleFunc wraps *mux.Router, allowing other packages to register with the router.
func (srv *Server) HandleFunc(path string, f func(http.ResponseWriter, *http.Request), methods ...string) {
	if len(methods) == 0 {
		methods = []string{http.MethodGet}
	}

	srv.r.HandleFunc(path, f).Methods(methods...)
}

// Fail renders an error template. The default behavior is to render 500.tmpl with
// an http.StatusInternalServerError status code. Fail also checks for application specific failures
// like form validation.
// The first keyvals value could be the name of a template to use instead of 500.tmpl.
func (srv *Server) Fail(ctx context.Context, w http.ResponseWriter, err error, keyvals ...interface{}) {
	tpl := "500.tmpl"
	if len(keyvals) > 0 && strings.HasSuffix(keyvals[0].(string), ".tmpl") {
		tpl, keyvals = keyvals[0].(string), keyvals[1:]
	}

	var validationErr interface {
		error
		Invalid() map[string]string
	}

	if errors.As(err, &validationErr) {
		srv.RenderTemplate(ctx, w, tpl, Data{}.
			FormErrors(validationErr.Invalid()).
			WithLog(err, keyvals...).
			WithCode(http.StatusBadRequest),
		)
		return
	}

	srv.RenderTemplate(ctx, w, tpl, Data{}.
		WithLog(err, keyvals...).
		WithCode(http.StatusInternalServerError),
	)
}

// RenderTemplate renders HTML templates.
func (srv *Server) RenderTemplate(ctx context.Context, w http.ResponseWriter, name string, data Data) {
	logger := log.FromContext(ctx)

	data["trace_id"] = log.TraceID(ctx)
	data["siteName"] = srv.siteName

	if logErr, ok := data[dLogErr]; ok {
		var kv []interface{}
		kv = append(kv, "err", logErr)
		extras, ok := data[dLogKV]
		if ok {
			kv = append(kv, extras.([]interface{})...)
		}

		// log the template name, avoiding loops to srv.Fail
		lvl := log.Info
		if name != "500.tmpl" {
			kv = append(kv, "template", name)
			lvl = log.Debug
		}

		lvl(logger).Log(kv...)
	}

	tmpl, ok := srv.templates[name]
	if !ok {
		srv.Fail(ctx, w, errors.New("no such template"), "template", name)
		return
	}

	// create a buffer to call ExecuteTemplate with, allowing for extra error handling
	// TODO: benchmark for allocations here
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "base.tmpl", data); err != nil && name != "500.tmpl" {
		srv.Fail(ctx, w, err, "msg", "executing template", "template", name)
		return
	} else if err != nil {
		log.Info(logger).Log("msg", "500 template failed to render", "err", err)
		return
	}

	if code, ok := data[dHTTPCode]; ok {
		w.WriteHeader(code.(int))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	buf.WriteTo(w)
	log.Debug(logger).Log("msg", "rendered template", "template", name)
}

// loadTemplates loads the UI from disk and caches it in a map.
// A lot of the inspiration came from an article I came across:
// https://blog.questionable.services/article/approximating-html-template-inheritance/
// Changes/imporvements to consider:
//   - make the layouts/includes locations configurable.
//   - allow loading from object storage (gcs/s3) instead of a local disk.
//   - support reloading with SIGHUP/other listeners.
//     Today, SIGHUP reloads the entire process, which works okay...
func (srv *Server) loadTemplates() error {
	layouts, err := filepath.Glob("ui/layouts/*.tmpl")
	if err != nil {
		return fmt.Errorf("load layouts: %s", err)
	}

	includes, err := filepath.Glob("ui/includes/*.tmpl")
	if err != nil {
		return fmt.Errorf("load includes: %s", err)
	}

	srv.mu.Lock()
	defer srv.mu.Unlock()

	for _, tpl := range includes {
		files := append(layouts, tpl)
		srv.templates[filepath.Base(tpl)] = template.Must(
			template.New(filepath.Base(tpl)).ParseFiles(files...),
		)
	}

	return nil
}

func (srv *Server) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				srv.Fail(r.Context(), w, fmt.Errorf("panic: %v", err), "msg", "recover panic", "debug_stack", string(debug.Stack()))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func csrfDisabled(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.FromContext(r.Context())
		log.Info(logger).Log("msg", "CSRF Protection disabled", "reason", "CSRF key not set.")
		next.ServeHTTP(w, r)
	})
}

func (srv *Server) csrf(next http.Handler) http.Handler {
	mw := csrfDisabled
	if string(srv.csrfKey) != "" {
		mw = csrf.Protect(
			srv.csrfKey,
			csrf.CookieName(srv.csrfCookieName),
			csrf.FieldName(srv.csrfFieldName),
			csrf.ErrorHandler(http.HandlerFunc(srv.csrfErrorHandler)),
		)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Excluding assets from setting the cookie.
		// Refactor into a switch or callback? if adding other paths.
		if strings.HasPrefix(r.URL.Path, "/assets/") {
			next.ServeHTTP(w, r)
			return
		}

		mw(next).ServeHTTP(w, r)
	})
}

func (srv *Server) csrfErrorHandler(w http.ResponseWriter, r *http.Request) {
	switch err := csrf.FailureReason(r); {
	// TODO: add 403 cases
	default:
		srv.Fail(r.Context(), w, err, "msg", "csrf.Protect encountered an error")
		return
	}
}

func (srv *Server) notFound(w http.ResponseWriter, r *http.Request) {
	srv.RenderTemplate(r.Context(), w, "404.tmpl", Data{}.WithCode(http.StatusNotFound))
}

func (srv *Server) indexPage(w http.ResponseWriter, r *http.Request) {
	srv.RenderTemplate(r.Context(), w, "home.tmpl", Data{})
}

// SESSION

var pubicPagePrefix = []string{
	"/login",
	"/assets/",
	"/forgot",
	"/register",
}

func isPublic(path string) bool {
	for _, k := range pubicPagePrefix {
		if strings.HasPrefix(path, k) {
			return true
		}
	}
	return false
}

func (srv *Server) sessionFromRequest(r *http.Request) (*session.Session, error) {
	ctx := r.Context()

	cookie, err := r.Cookie(srv.authCookieName)
	if err != nil {
		return nil, err
	}

	value := make(map[string]string)
	if err := srv.cookie.Decode(srv.authCookieName, cookie.Value, &value); err != nil {
		return nil, err
	}

	id, ok := value["id"]
	if !ok {
		return nil, errors.New("auth cookie present but no id value")
	}

	sess, err := srv.sessiondb.FindSession(ctx, id)
	if err != nil {
		return nil, err
	}

	if sess.CreatedAt.Before(time.Now().UTC().Add(time.Duration(-30) * time.Minute)) {
		return nil, fmt.Errorf("session (id %s) expired", sess.ID)
	}

	return sess, nil
}

func (srv *Server) authMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx    = r.Context()
			logger = log.FromContext(ctx)
		)
		sess, err := srv.sessionFromRequest(r)
		if err == http.ErrNoCookie && isPublic(r.URL.Path) {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		} else if err != nil && isPublic(r.URL.Path) {
			srv.Fail(r.Context(), w, err, "msg", "auth mw failed to render")
			return
		} else if err != nil {
			// TODO: handle not found and destroy the session somewhere.
			log.Debug(logger).Log("err", err, "msg", "check auth")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		ctx = viewer.NewContext(ctx, viewer.Viewer{
			UserID:    sess.UserID,
			SessionID: sess.ID,
		})

		if strings.HasPrefix(r.URL.Path, "/login") {
			level.Debug(logger).Log("msg", "redirect session to dashboard page", "session_id", sess.ID)
			http.Redirect(w, r, "/dashboard", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
