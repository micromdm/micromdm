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
	"sync"
	"text/template"

	"github.com/gorilla/mux"

	"micromdm.io/v2/pkg/log"
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

// Server implements Framework.
type Server struct {
	r *mux.Router

	mu        sync.Mutex
	templates map[string]*template.Template

	siteName string
}

// Config parameters to create a new Server.
type Config struct {
	Logger   log.Logger
	SiteName string
}

// New creates a Server.
func New(config Config) (*Server, error) {
	srv := &Server{
		r:         mux.NewRouter(),
		templates: make(map[string]*template.Template),
		siteName:  config.SiteName,
	}

	srv.r.Use(
		log.HTTP(config.Logger), // HTTP logging middleware.
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

// HandleFunc wraps *mux.Router, allowing other packages to register with the router.
func (srv *Server) HandleFunc(path string, f func(http.ResponseWriter, *http.Request), methods ...string) {
	if len(methods) == 0 {
		methods = []string{http.MethodGet}
	}

	srv.r.HandleFunc(path, f).Methods(methods...)
}

// Fail renders the 500 InternalServerError template and logs accordingly.
func (srv *Server) Fail(ctx context.Context, w http.ResponseWriter, err error, keyvals ...interface{}) {
	srv.RenderTemplate(ctx, w, "500.tmpl", Data{}.
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
		if name != "500.tmpl" {
			kv = append(kv, "template", name)
		}

		log.Info(logger).Log(kv...)
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
				srv.Fail(r.Context(), w, fmt.Errorf("panic: %v", err), "msg", "recover panic")
				debug.PrintStack()
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (srv *Server) notFound(w http.ResponseWriter, r *http.Request) {
	srv.RenderTemplate(r.Context(), w, "404.tmpl", Data{}.WithCode(http.StatusNotFound))
}

func (srv *Server) indexPage(w http.ResponseWriter, r *http.Request) {
	srv.RenderTemplate(r.Context(), w, "home.tmpl", Data{})
}
