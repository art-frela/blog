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
	bf "gopkg.in/russross/blackfriday.v2"
)

const (
	templatePATH = "./assets/templates/*.html"
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

// RedirectToPosts - simple redirect for posts url
func (pc *PostController) RedirectToPosts(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
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
	// if len(posts) == 0 {
	// 	err = fmt.Errorf("not found no one post in the repository")
	// 	render.Render(w, r, ErrNotFound(err))
	// 	return
	// }
	for i, p := range posts {
		posts[i].Content = template.HTML(bf.Run([]byte(p.Content)))
	}
	data := templatePostsFill{
		Title: "POSTS",
		Posts: posts,
	}
	tmpl := template.Must(template.New("indexPOST").ParseGlob(templatePATH))
	tmpl.ExecuteTemplate(w, "indexPOST", data)
}

// GetOnePost returns the one specified by id post from storage
func (pc *PostController) GetOnePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	post, err := pc.PostRepo.FindByID(id)
	if err != nil {
		render.Render(w, r, ErrServerInternal(err))
		return
	}
	post.Content = template.HTML(bf.Run([]byte(post.Content))) // use blackfriday for markdown to html view
	data := templateOnePostFill{
		Title: post.Title,
		Post:  post,
	}
	tmpl := template.Must(template.New("indexSinglePOST").ParseGlob(templatePATH))
	tmpl.ExecuteTemplate(w, "indexSinglePOST", data)
}

// EditPost - handler func for expose edit form for Posts
func (pc *PostController) EditPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	post, err := pc.PostRepo.FindByID(id)
	if err != nil {
		render.Render(w, r, ErrServerInternal(err))
		return
	}
	data := templateOnePostFill{
		Title: post.Title,
		Post:  post,
	}
	tmpl := template.Must(template.New("indexEditPOST").ParseGlob(templatePATH))
	tmpl.ExecuteTemplate(w, "indexEditPOST", data)
}

// WriteNewPost - handler func for expose edit form for new Post
func (pc *PostController) WriteNewPost(w http.ResponseWriter, r *http.Request) {
	post := domain.PostInBlog{
		Title:   "",
		Content: "",
	}
	data := templateOnePostFill{
		Title: post.Title,
		Post:  post,
	}
	tmpl := template.Must(template.New("indexNewPOST").ParseGlob(templatePATH))
	tmpl.ExecuteTemplate(w, "indexNewPOST", data)
}

// UpdPost - handler func for update post in the Storage
func (pc *PostController) UpdPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	params := &NewPostRequest{}
	if err := render.Bind(r, params); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	newpost := domain.PostInBlog{
		ID:      id,
		Title:   params.Title,
		Content: template.HTML(params.Content),
		Rubric: domain.Rubric{
			ID: params.RubricID,
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
	render.Render(w, r, OkStatus(id))
}

// AddNewPost - handler func for save new post in the storage
func (pc *PostController) AddNewPost(w http.ResponseWriter, r *http.Request) {
	params := &NewPostRequest{}
	if err := render.Bind(r, params); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	// contentBts := []byte(params.Content)
	// contentMD := bf.Run(contentBts)
	// contentSafeHTML := bluemonday.UGCPolicy().SanitizeBytes(contentMD)
	newpost := domain.PostInBlog{
		Title:   params.Title,
		Content: template.HTML(params.Content),
		Rubric: domain.Rubric{
			ID: params.RubricID,
		},
	}
	id, err := pc.PostRepo.Save(newpost)
	if err != nil {
		err = fmt.Errorf("try to save new post %v, error %v", newpost, err)
		render.Render(w, r, ErrServerInternal(err))
		return
	}
	render.Render(w, r, OkStatusCreated(id))
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

// NewPostRequest contract with front-end for posts creating
type NewPostRequest struct {
	Title    string `json:"title"`
	RubricID string `json:"rubric_id"`
	Content  string `json:"content"`
	UserID   string `json:"user_id"`
}

// Bind - implement Bind method for chi.render interface
func (npr *NewPostRequest) Bind(r *http.Request) error {
	return nil
}

// SuccessResponse structure for json response success results
type SuccessResponse struct {
	Message        string `json:"message"`  // text of message
	HTTPStatusCode int    `json:"httpcode"` // http response status code
	StatusText     string `json:"status"`   // user-level status message
}

// Render - implement Render method for chi.render interface
func (sr *SuccessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, sr.HTTPStatusCode)
	return nil
}

// OkStatusCreated rendered func for HTTP 201 Created response
func OkStatusCreated(message string) render.Renderer {
	return &SuccessResponse{
		Message:        message,
		HTTPStatusCode: http.StatusCreated,
		StatusText:     http.StatusText(http.StatusCreated),
	}
}

// OkStatus rendered func for HTTP 200 OK response
func OkStatus(message string) render.Renderer {
	return &SuccessResponse{
		Message:        message,
		HTTPStatusCode: http.StatusOK,
		StatusText:     http.StatusText(http.StatusOK),
	}
}
