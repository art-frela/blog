package infra

import (
	"fmt"

	"github.com/art-frela/blog/domain"
	"github.com/gofrs/uuid"
)

// SimplePostRepo implement in-memory storage for PostRepository
type SimplePostRepo struct {
	Posts map[string]domain.PostInBlog
	//Comments map[string]domain.CommentOfPost
}

// NewSimplePostRepo return SimplePostRepo with filled example posts
func NewSimplePostRepo(n int) *SimplePostRepo {
	spr := &SimplePostRepo{}
	return spr.fillExampleData(n)
}

// FindByID implement post repository for map[string]Posts
func (spr *SimplePostRepo) FindByID(id string) (domain.PostInBlog, error) {
	post, ok := spr.Posts[id]
	if !ok {
		err := fmt.Errorf("post with id=%s not found", id)
		return post, err
	}
	return post, nil
}

// Find implement post repository for map[string]Posts
// returns slice of posts from map[string]Posts
func (spr *SimplePostRepo) Find(limit, offset int) ([]domain.PostInBlog, error) {
	// !!!UNSORTED MAP STORAGE!!! TODO: use DB with ordering
	posts := make([]domain.PostInBlog, 0, limit)
	if len(spr.Posts) == 0 {
		err := fmt.Errorf("no one post in the storage, you can be first ;-)")
		return posts, err
	}
	//
	inLimit, inOffset, ix := 0, 0, 0
	for _, p := range spr.Posts {
		if ix >= inOffset {
			if inLimit < limit {
				inLimit++
				posts = append(posts, p)
			}
			ix++
		}
	}
	return posts, nil
}

// Save implement post repository for map[string]Posts
// add new post to the map
func (spr *SimplePostRepo) Save(p domain.PostInBlog) (string, error) {
	newID := uuid.Must(uuid.NewV4()).String()
	p.ID = newID
	spr.Posts[newID] = p
	return newID, nil
}

// Update implement post repository for map[string]Posts
// update exists post in the map
func (spr *SimplePostRepo) Update(p domain.PostInBlog) error {
	_, ok := spr.Posts[p.ID]
	if !ok {
		err := fmt.Errorf("post with id=%s not found, your update is rejected", p.ID)
		return err
	}
	spr.Posts[p.ID] = p
	return nil
}

// fillExampleData fills SimplePostRepo with fake posts exactly N pieces,
// but no more 50th
func (spr *SimplePostRepo) fillExampleData(n int) *SimplePostRepo {
	if n > 50 || n <= 0 { // simple fuse
		n = 10
	}
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
	if spr.Posts == nil {
		spr.Posts = make(map[string]domain.PostInBlog)
	}

	for i := 1; i <= n; i++ {
		post := domain.PostInBlog{
			Title:   fmt.Sprintf(postTmpl.Title, i),
			Author:  postTmpl.Author,
			Rubric:  postTmpl.Rubric,
			Content: postTmpl.Content,
		}
		post.Author.ID = fmt.Sprintf(postTmpl.Author.ID, i)
		post.Author.Name = fmt.Sprintf(postTmpl.Author.Name, i)
		post.Rubric.ID = fmt.Sprintf(postTmpl.Rubric.ID, i)
		newID := uuid.Must(uuid.NewV4()).String()
		post.SetID(newID)
		spr.Posts[newID] = post
	}
	return spr
}
