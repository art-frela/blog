package infra

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/art-frela/blog/domain"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

const (
	templatePOSTS = "./assets/*.html"
	postID        = "id"
	//httpTimeOut   = 30 * time.Second
)

// [HANDLER FUNCS]

// PostController - main controller for Posts
type PostController struct {
	PostRepo domain.PostRepository
	//CommentsRepo domain.CommentsRepository
}

// NewPostController is a builder for PostController
func NewPostController(repo domain.PostRepository) *PostController {
	pc := &PostController{
		PostRepo: repo,
	}
	return pc
}

// GetPosts - handler func for search query text at the Sites
func (pc *PostController) GetPosts(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	offset, _ := strconv.Atoi(r.FormValue("offset"))
	if limit == 0 {
		limit = 50
	}
	posts, err := pc.PostRepo.Find(limit, offset)
	if err != nil {
		render.Render(w, r, ErrServerInternal(err))
		return
	}
	if len(posts) == 0 {
		err = fmt.Errorf("not found no one post in the repository")
		render.Render(w, r, ErrNotFound(err))
		return
	}
	data := templatePostsFill{
		Title: "POSTS",
		Posts: posts,
	}
	tmpl := template.Must(template.New("indexPOST").ParseGlob(templatePOSTS))
	tmpl.ExecuteTemplate(w, "indexPOST", data)
}

// GetOnePost - handler func for search query text at the Sites
func (pc *PostController) GetOnePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, postID)
	post, err := pc.PostRepo.FindByID(id)
	if err != nil {
		render.Render(w, r, ErrServerInternal(err))
		return
	}
	data := templateOnePostFill{
		Title: post.Title,
		Post:  post,
	}
	tmpl := template.Must(template.New("indexSinglePOST").ParseGlob(templatePOSTS))
	tmpl.ExecuteTemplate(w, "indexSinglePOST", data)
}

// EditPost - handler func for exposeedit form for Posts
func (pc *PostController) EditPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, postID)
	post, err := pc.PostRepo.FindByID(id)
	if err != nil {
		render.Render(w, r, ErrServerInternal(err))
		return
	}
	data := templateOnePostFill{
		Title: post.Title,
		Post:  post,
	}
	tmpl := template.Must(template.New("indexEditPOST").ParseGlob(templatePOSTS))
	tmpl.ExecuteTemplate(w, "indexEditPOST", data)
}

// WriteNewPost - handler func for expose edit form for new Posts
func (pc *PostController) WriteNewPost(w http.ResponseWriter, r *http.Request) {
	post := domain.PostInBlog{
		Title:   "",
		Content: "",
	}
	data := templateOnePostFill{
		Title: post.Title,
		Post:  post,
	}
	tmpl := template.Must(template.New("indexNewPOST").ParseGlob(templatePOSTS))
	tmpl.ExecuteTemplate(w, "indexNewPOST", data)
}

// UpdPost - handler func for update post in the Storage
func (pc *PostController) UpdPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, postID)
	newpost := domain.PostInBlog{
		ID:      id,
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
		Rubric: domain.Rubric{
			Title: r.FormValue("rubric"),
		},
	}
	oldpost, err := pc.PostRepo.FindByID(id)
	if err != nil {
		render.Render(w, r, ErrServerInternal(err))
		return
	}
	// Simple comparison and fill values for upd Post
	// TODO: add comparison/merge method for PostInBlog in the domain.go, without reflection please!!!
	if oldpost.Title != newpost.Title {
		oldpost.Title = newpost.Title
	}
	if oldpost.Content != newpost.Content {
		oldpost.Content = newpost.Content
	}
	oldpost.ModifiedAt = time.Now().Format(time.RFC3339)
	if oldpost.Rubric.Title != newpost.Rubric.Title {
		oldpost.Rubric.Title = newpost.Rubric.Title
	}

	err = pc.PostRepo.Update(oldpost)
	if err != nil {
		render.Render(w, r, ErrServerInternal(err))
		return
	}
	w.WriteHeader(http.StatusOK)
}

// AddNewPost - handler func for save new post in the storage
func (pc *PostController) AddNewPost(w http.ResponseWriter, r *http.Request) {
	newpost := domain.PostInBlog{
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
		Rubric: domain.Rubric{
			Title: r.FormValue("rubric"),
		},
	}
	id, err := pc.PostRepo.Save(newpost)
	if err != nil {
		render.Render(w, r, ErrServerInternal(err))
		return
	}
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("/posts/" + id))
}

type templatePostsFill struct {
	Title string
	Posts []domain.PostInBlog
}

type templateOnePostFill struct {
	Title string
	Post  domain.PostInBlog
}

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render - implement method Render for render.Renderer
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest - wrapper for make err structure
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

// ErrServerInternal - wrapper for make err structure
func ErrServerInternal(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal server error.",
		ErrorText:      err.Error(),
	}
}

// ErrNotFound - wrapper for make err structure for empty result
func ErrNotFound(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     http.StatusText(http.StatusNotFound),
		ErrorText:      err.Error(),
	}
}

// ErrUnsupportedFormat - 415 error implementation
var ErrUnsupportedFormat = &ErrResponse{HTTPStatusCode: http.StatusUnsupportedMediaType, StatusText: "415 - Unsupported Media Type. Please send JSON"}
