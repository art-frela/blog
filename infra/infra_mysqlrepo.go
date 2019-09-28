package infra

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/art-frela/blog/domain"
	"github.com/art-frela/blog/models"
	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/gofrs/uuid"
)

// MySQLPostRepository - post repository implementation
type MySQLPostRepository struct {
	db  *sql.DB
	log *logrus.Logger
	ctx context.Context
}

// NewMySQLPostRepository returns MySQL post repository
func NewMySQLPostRepository(mysqlURL string, logger *logrus.Logger, countExamplePosts int, clearStorage bool) *MySQLPostRepository {
	repo := &MySQLPostRepository{}
	repo.log = logger
	db, err := sql.Open("mysql", mysqlURL)
	if err != nil {
		repo.log.Fatalf("error open mysql, %v", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	err = db.Ping()
	if err != nil {
		repo.log.Fatalf("error open connection to mysql server, %v", err)
	}
	repo.db = db
	boil.SetDB(db)
	ctx := context.Background()
	repo.ctx = ctx
	if clearStorage {

		rowsAff, err := models.Posts().DeleteAll(repo.ctx, repo.db)
		if err != nil {
			repo.log.Errorf("error delete all posts from database, %v", err)
		}
		repo.log.Infof("deleted all (%d) posts from database, ", rowsAff)

	}
	repo.fillExampleData(countExamplePosts)
	return repo
}

// FindByID implement post repository for map[string]Posts
func (myr *MySQLPostRepository) FindByID(id string) (domain.PostInBlog, error) {
	post := domain.PostInBlog{}
	modelPost, err := models.FindPost(myr.ctx, myr.db, id)
	if err != nil {
		return post, err
	}
	post = convertModelPostToDomainPost(*modelPost)
	return post, nil
}

// Find implement post repository for mysql
// returns slice of posts
func (myr *MySQLPostRepository) Find(limit, offset int) ([]domain.PostInBlog, error) {
	posts := make([]domain.PostInBlog, 0, 16)
	modelPosts, err := models.Posts(qm.Limit(limit), qm.Offset(offset)).All(myr.ctx, myr.db)
	if err != nil {
		return posts, err
	}
	for _, p := range modelPosts {
		posts = append(posts, convertModelPostToDomainPost(*p))
	}
	return posts, nil
}

// Save implement post repository for MySQL
// add new post to the DB
func (myr *MySQLPostRepository) Save(p domain.PostInBlog) (string, error) {
	newID := uuid.Must(uuid.NewV4()).String()
	templPost := p.GetTemplatePost()
	p.ID = newID
	p.State = templPost.State
	// TODO: after implement rubric and author at the fron GUI, del this mock
	p.Rubric.ID = "00000000-0000-0000-00000000"
	p.Author.ID = domain.AnonimousID
	modelPost := convertDomainPostToModelPost(p)
	err := modelPost.Insert(myr.ctx, myr.db, boil.Infer())
	if err != nil {
		return "", err
	}
	return newID, nil
}

// Update implement post repository for map[string]Posts
// update exists post in the map
func (myr *MySQLPostRepository) Update(p domain.PostInBlog) error {
	// TODO: add validator fot id, title, content etc...
	postModel, err := models.FindPost(myr.ctx, myr.db, p.ID)
	if err != nil {
		return err
	}
	*postModel = convertDomainPostToModelPost(p)
	_, err = postModel.Update(myr.ctx, myr.db, boil.Infer())
	return err
}

// fillExampleData fills SimplePostRepo with fake posts exactly N pieces,
// but no more 3rd
func (myr *MySQLPostRepository) fillExampleData(n int) {
	if n > 3 || n <= 0 { // simple fuse
		n = 3
	}
	//newID := uuid.Must(uuid.NewV4()).String()
	postTmpl := domain.PostInBlog{
		Title: "Example post #%d",
		Author: domain.User{
			ID:    "someUserID_%d",
			Name:  "anonimous#%d",
			EMail: "anonim@example.com",
		},
		Rubric: domain.Rubric{
			ID:    "rubricID_%d",
			Title: "Go for fun",
		},
		Content: `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore
            magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
            consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
            pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est
			laborum.`,
	}

	for i := 1; i <= n; i++ {
		post := domain.PostInBlog{
			Title:   fmt.Sprintf(postTmpl.Title, i),
			Author:  postTmpl.Author,
			Rubric:  postTmpl.Rubric,
			Content: postTmpl.Content,
		}
		post.Author.ID = domain.AnonimousID
		post.Author.Name = fmt.Sprintf(postTmpl.Author.Name, i)
		post.Rubric.ID = "00000000-0000-0000-00000000"
		newID := uuid.Must(uuid.NewV4()).String()
		post.SetID(newID)
		_, err := myr.Save(post)
		if err != nil {
			myr.log.Errorf("for post=%+v, error %v", post, err)
		}
	}
}

// convertModelPostToDomainPost - return domainPost make from model post
func convertModelPostToDomainPost(post models.Post) domain.PostInBlog {
	targetPost := domain.PostInBlog{}
	targetPost.SetID(post.ID)
	targetPost.SetContent(post.Content)
	targetPost.Author.ID = post.AuthorID.String
	targetPost.SetTitle(post.Title)
	targetPost.Rubric.ID = post.RubricID.String
	return targetPost
}

// convertDomainPostToModelPost - return model post  make from domain post
func convertDomainPostToModelPost(post domain.PostInBlog) models.Post {
	targetPost := models.Post{}
	targetPost.ID = post.ID
	targetPost.Title = post.Title
	targetPost.Content = string(post.Content)
	targetPost.AuthorID.String = post.Author.ID
	targetPost.RubricID.String = post.Rubric.ID
	targetPost.State.String = post.State
	targetPost.CountOfViews = int(post.CountOfViews)
	return targetPost
}
