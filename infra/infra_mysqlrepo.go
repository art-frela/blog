package infra

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/art-frela/blog/domain"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofrs/uuid"
)

// MySQLPostRepository - post repository implementation
type MySQLPostRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

// NewMySQLPostRepository returns MySQL post repository
func NewMySQLPostRepository(mysqlURL string, logger *logrus.Logger) *MySQLPostRepository {
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
	repo.fillExampleData(20)
	return repo
}

// FindByID implement post repository for map[string]Posts
func (myr *MySQLPostRepository) FindByID(id string) (domain.PostInBlog, error) {
	q := `select 
			id, title, author_id, rubric_id, COALESCE(tags, ''), state, content, created_at, modified_at, 
			COALESCE(parent_post_id, ''), count_of_views, count_of_stars, COALESCE(comments_ids, '') 
		  from posts 
		  where id=?;`
	row := myr.db.QueryRow(q, id)
	var post domain.PostInBlog
	var author domain.User
	var tags, comments string // string
	err := row.Scan(&post.ID, &post.Title, &author.ID, &post.Rubric.ID, &tags, &post.State, &post.Content,
		&post.CreatedAt, &post.ModifiedAt, &post.ParentPostID, &post.CountOfViews,
		&post.CountOfStars, &comments)
	if err != nil {
		return post, err
	}

	//TODO: add search user by ID and fill user structure
	post.SetAuthor(author)
	postTags := &domain.Tags{}
	err = json.NewDecoder(strings.NewReader(string(tags))).Decode(postTags)
	if err != nil {
		myr.log.Warnf("error unmarshalling tags, %v, set default value [\"post\"]", err)
		*postTags = domain.Tags{"post"}
	}
	post.SetTags(*postTags)
	// TODO: set task to get comments

	return post, nil
}

// Find implement post repository for mysql
// returns slice of posts
func (myr *MySQLPostRepository) Find(limit, offset int) ([]domain.PostInBlog, error) {
	posts := make([]domain.PostInBlog, 0, 16)
	q := fmt.Sprintf(`select 
			id, title, author_id, rubric_id, COALESCE(tags, ''), state, content, created_at, modified_at, 
			COALESCE(parent_post_id, ''), count_of_views, count_of_stars, COALESCE(comments_ids, '') 
		  from posts where state='%s' order by count_of_stars desc limit ? offset ?;`,
		domain.PostStatePublic)
	rows, err := myr.db.Query(q, limit, offset)
	if err != nil {
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		var post domain.PostInBlog
		var author domain.User
		var tags, comments string
		err := rows.Scan(&post.ID, &post.Title, &author.ID, &post.Rubric.ID, &tags, &post.State, &post.Content,
			&post.CreatedAt, &post.ModifiedAt, &post.ParentPostID, &post.CountOfViews,
			&post.CountOfStars, &comments)
		if err != nil {
			return posts, err
		}
		//TODO: add search user by ID and fill user structure
		post.SetAuthor(author)
		postTags := &domain.Tags{}
		err = json.NewDecoder(strings.NewReader(tags)).Decode(postTags)
		if err != nil {
			myr.log.Warnf("error unmarshalling tags, %v, set default value [\"post\"]", err)
			*postTags = domain.Tags{"post"}
		}
		post.SetTags(*postTags)
		// TODO: set task to get comments
		posts = append(posts, post)
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
	// TODO: after implement rubric and author at the fron GUI, del thih mock
	p.Rubric.ID = "00000000-0000-0000-00000000"
	p.Author.ID = domain.AnonimousID
	q := `insert into posts (id, title, rubric_id, content, author_id, state) values(?,?,?,?,?,?)`
	result, err := myr.db.Exec(q, p.ID, p.Title, p.Rubric.ID, p.Content, p.Author.ID, p.State)
	if err != nil {
		return newID, err
	}
	rowAffected, _ := result.RowsAffected()
	myr.log.Debug("new post create succesfull, row affected %d", rowAffected)
	return newID, nil
}

// Update implement post repository for map[string]Posts
// update exists post in the map
func (myr *MySQLPostRepository) Update(p domain.PostInBlog) error {
	// TODO: add validator fot id, title, content etc...
	q := `update posts set title=?, content=? where id=?`
	result, err := myr.db.Exec(q, p.Title, p.Content, p.ID)
	if err != nil {
		return err
	}
	rowAffected, _ := result.RowsAffected()
	myr.log.Debugf("update postID=%s (%d rows affected)", p.ID, rowAffected)
	return nil
}

// fillExampleData fills SimplePostRepo with fake posts exactly N pieces,
// but no more 50th
func (myr *MySQLPostRepository) fillExampleData(n int) {
	if n > 50 || n <= 0 { // simple fuse
		n = 10
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
