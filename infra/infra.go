package infra

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

// BlogServer -
type BlogServer struct {
	log        *logrus.Logger
	mux        *chi.Mux
	controller *PostController
}

// NewBlogServer is builder for BlogServer
func NewBlogServer(mysqlURL string) *BlogServer {
	bs := &BlogServer{}
	logger := logrus.New()
	pr := NewMySQLPostRepository(mysqlURL, logger)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(customHTTPLogger)
	bs.mux = r
	bs.log = logger
	bs.controller = NewPostController(pr)
	return bs
}

// Run is running blogServer
func (bs *BlogServer) Run(hostPort string) {
	bs.registerPostRoutes()
	bs.log.Infof("http server starting on the [%s] tcp port", hostPort)
	bs.log.Fatal(http.ListenAndServe(hostPort, bs.mux))
}

func (bs *BlogServer) registerPostRoutes() {
	bs.mux.Route("/posts", func(r chi.Router) {
		r.Get("/", bs.controller.GetPosts)
		r.Get("/{"+postID+"}", bs.controller.GetOnePost)
		r.Get("/{"+postID+"}/edit", bs.controller.EditPost)
		r.Get("/new", bs.controller.WriteNewPost)
		r.Put("/{"+postID+"}", bs.controller.UpdPost)
		r.Post("/", bs.controller.AddNewPost)
	})
}

// [CUSTOM MIDDLEWARE]

// filterContentType - middleware to check content type as JSON
func filterContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Filtering requests by MIME type
		if r.Method == "POST" { // filter for POST request
			if r.Header.Get("Content-type") != "application/json" {
				render.Render(w, r, ErrUnsupportedFormat)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// customHTTPLogger - middleware to logrus logger
func customHTTPLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start).String()
		host, _ := os.Hostname()
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"proto":  r.Proto,
			"remote": r.RemoteAddr,
			"url":    r.RequestURI,
			//"code":     r.Response.StatusCode,
			"duration": duration,
		}).Infof("%s", host)
	})
}
