package infra

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/art-frela/blog/domain"

	_ "github.com/art-frela/blog/docs"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	swag "github.com/swaggo/http-swagger"
)

// BlogServer -
type BlogServer struct {
	log        *logrus.Logger
	mux        *chi.Mux
	controller *PostController
}

// NewBlogServer is builder for BlogServer
func NewBlogServer(storageURL string, countExamplePosts int, clearStorage bool) *BlogServer {
	bs := &BlogServer{}
	logger := logrus.New()
	storageType := defineStorageType(storageURL)
	pr := NewPostStorage(storageType, storageURL, logger, countExamplePosts, clearStorage)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(customHTTPLogger)
	// add aka fileserver
	filesDir := filepath.Join(".", "assets/css")
	FileServer(r, "/css", http.Dir(filesDir))
	filesDir = filepath.Join(".", "assets/js")
	FileServer(r, "/js", http.Dir(filesDir))
	filesDir = filepath.Join(".", "assets/img")
	FileServer(r, "/img", http.Dir(filesDir))
	bs.mux = r
	bs.log = logger
	bs.controller = NewPostController(pr)
	return bs
}

// NewPostStorage looks like AbstractFactory of PostRepositories
func NewPostStorage(storageType string, storageURL string, logger *logrus.Logger, countExamplePosts int, clearStorage bool) domain.PostRepository {
	switch storageType {
	case "mysql":
		return NewMySQLPostRepository(storageURL, logger, countExamplePosts, clearStorage)
	default:
		return NewMongoPostRepo(storageURL, logger, countExamplePosts, clearStorage)
	}
}

// Run is running blogServer
func (bs *BlogServer) Run(hostPort string) *http.Server {
	srv := &http.Server{Addr: hostPort, Handler: bs.mux}
	bs.registerRoutes(hostPort)
	bs.log.Infof("http server starting on the [%s] tcp port", hostPort)
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			bs.log.Fatalf("http server error: %v", err)
		}
	}()
	return srv
}

// Stop is stopping blogServer
func (bs *BlogServer) Stop(srv *http.Server) {
	bs.log.Info("http server stopping")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		bs.log.Errorf("http server stopping error, %v", err)
	}
}

func (bs *BlogServer) registerRoutes(hostPort string) {
	uri := fmt.Sprintf("http://%s/swagger/doc.json", hostPort)
	bs.mux.Get("/swagger/*", swag.Handler(
		swag.URL(uri), //The url pointing to API definition"
	))
	bs.mux.Route("/posts", func(r chi.Router) {
		r.Get("/", bs.controller.GetPosts)
		r.Get("/{id}", bs.controller.GetOnePost)
		r.Get("/{id}/edit", bs.controller.EditPost)
		r.Get("/new", bs.controller.WriteNewPost)

		//r.Post("/", bs.controller.AddNewPost)
	})
	bs.mux.Route("/api/v1", func(r chi.Router) {
		r.Route("/posts", func(r chi.Router) {
			r.Use(filterContentType)
			//r.Get("/", bs.controller.GetPostJSON) // TODO: implement
			r.Post("/", bs.controller.AddNewPost)
			r.Put("/{id}", bs.controller.UpdPost)
		})
	})
	bs.mux.Route("/", func(r chi.Router) {
		r.Get("/", bs.controller.RedirectToPosts)
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

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

// defineStorageType helper for define storage type by connection string
func defineStorageType(storageURL string) string {
	storageType := "mysql"
	if strings.HasPrefix(strings.ToLower(storageURL), "mongodb") {
		storageType = "mongo"
	}
	return storageType
}
