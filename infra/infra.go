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

	_ "github.com/art-frela/blog/docs" // docs is swagger generated file, don't modified it!
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	swag "github.com/swaggo/http-swagger"
)

// ContextID is our type to retrieve our context
// objects
type contextStatusID int

const (
	StatusCtxKey contextStatusID = 0
)

// BlogServer -
type BlogServer struct {
	log        *logrus.Entry
	mux        *chi.Mux
	controller *PostController
	config     *viper.Viper
	srv        *http.Server
}

// NewBlogServer is builder for BlogServer
func NewBlogServer(countExamplePosts int, clearStorage bool) *BlogServer {
	bs := &BlogServer{}
	bs.setConfig()
	bs.setLogger("0.0.2")
	storageType := defineStorageType(bs.config.GetString("database.url"))
	pr := NewPostStorage(storageType, bs.config.GetString("database.url"), bs.config.GetString("database.name"), bs.log, countExamplePosts, clearStorage)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	//r.Use(middleware.Logger)
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
	bs.controller = NewPostController(pr)
	return bs
}

// NewPostStorage looks like AbstractFactory of PostRepositories
func NewPostStorage(storageType string, storageURL, database string, logger *logrus.Entry, countExamplePosts int, clearStorage bool) domain.PostRepository {
	switch storageType {
	case "mysql":
		return NewMySQLPostRepository(storageURL, database, logger, countExamplePosts, clearStorage)
	default:
		return NewMongoPostRepo(storageURL, database, logger, countExamplePosts, clearStorage)
	}
}

// Run is running blogServer
func (bs *BlogServer) Run() {
	hostPort := fmt.Sprintf("%s:%s", bs.config.GetString("httpd.host"), bs.config.GetString("httpd.port"))
	srv := &http.Server{Addr: hostPort, Handler: bs.mux}
	bs.registerRoutes()
	bs.log.Infof("http server starting on the [%s] tcp port", hostPort)
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			bs.log.Fatalf("http server error: %v", err)
		}
	}()
	bs.srv = srv
}

// Stop is stopping blogServer
func (bs *BlogServer) Stop() {
	bs.log.Info("http server stopping")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := bs.srv.Shutdown(ctx); err != nil {
		bs.log.Errorf("http server stopping error, %v", err)
	}
}

func (bs *BlogServer) registerRoutes() {
	uri := fmt.Sprintf("http://%s:%s/swagger/doc.json", bs.config.GetString("swagger.host"), bs.config.GetString("httpd.port"))
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

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

// customHTTPLogger - middleware to logrus logger
func customHTTPLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cw := statusRecorder{w, http.StatusOK}
		start := time.Now()
		next.ServeHTTP(&cw, r)
		duration := time.Since(start).String()
		log := logrus.WithFields(logrus.Fields{
			"method":   r.Method,
			"proto":    r.Proto,
			"remote":   r.RemoteAddr,
			"url":      r.RequestURI,
			"code":     cw.status,
			"duration": duration,
		})

		host, _ := os.Hostname()
		switch {
		case cw.status < 300:
			log.Infof("%s", host)
		case cw.status > 300 && cw.status < 400:
			log.Warnf("%s", host)
		case cw.status >= 400:
			log.Errorf("%s", host)
		}
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
