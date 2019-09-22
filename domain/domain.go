package domain

const (
	// PostStateWrite - state of post when it only creating by Author
	PostStateWrite = "write"
	// PostStateModerate - is state of post when it's saved, but not showing for other users, only Author/Moderator/Admin cat access to it
	PostStateModerate = "moderate"
	// PostStatePublic - is state of post when it's saved and all can see it
	PostStatePublic = "public"
	// PostStateBlocked - is state of post when it's saved, but no one can't access to it, only admin
	PostStateBlocked = "blocked"

	// AnonimousID - default userID
	AnonimousID = "00000000-0000-0000-00000000"

	// UserAdmin - roleID for administrator
	UserAdmin = iota
	// UserModerator - roleID for moderator
	UserModerator
	// UserDefault - roleID for other users
	UserDefault
)

// PostInBlog - main entity my blog, like story in livejournal
type PostInBlog struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Author       User     `json:"author"`
	Rubric       Rubric   `json:"rubric"`
	Content      string   `json:"content"`
	Tags         Tags     `json:"tags"`
	State        string   `json:"state"`
	CreatedAt    string   `json:"created_at"`  // RFC3339/ISO8601
	ModifiedAt   string   `json:"modified_at"` // RFC3339/ISO8601
	ParentPostID string   `json:"parent_post_id"`
	CountOfViews int64    `json:"count_of_views"`
	CountOfStars int64    `json:"count_of_stars"`
	CommentsIDs  []string `json:"comments_ids"`
}

// PostRepository - storage of Posts
type PostRepository interface {
	FindByID(id string) (PostInBlog, error)
	Find(limit, offset int) ([]PostInBlog, error)
	//FindByRubric(r Rubric) ([]PostInBlog, error)
	//FindByQuery(phrase string) ([]PostInBlog, error)
	Save(p PostInBlog) (string, error)
	Update(p PostInBlog) error
	//DeletePost(p PostInBlog) (bool, error)
}

// [Setters fo PostInBlog]

// SetID - setter for ID
func (p *PostInBlog) SetID(id string) *PostInBlog {
	p.ID = id
	return p
}

// SetTitle - setter for Title
func (p *PostInBlog) SetTitle(title string) *PostInBlog {
	// TODO: add XSS checker
	p.Title = title
	return p
}

// SetAuthor - setter for Author
func (p *PostInBlog) SetAuthor(author User) *PostInBlog {
	p.Author = author
	return p
}

// SetRubric - setter for ID
func (p *PostInBlog) SetRubric(rubric Rubric) *PostInBlog {
	p.Rubric = rubric
	return p
}

// SetContent - setter for ID
func (p *PostInBlog) SetContent(content string) *PostInBlog {
	// TODO: add XSS checker
	p.Content = content
	return p
}

// SetStateWrite - setter for state to write
func (p *PostInBlog) SetStateWrite() *PostInBlog {
	p.State = PostStateWrite
	return p
}

// SetStateModerate - setter for state to write
func (p *PostInBlog) SetStateModerate() *PostInBlog {
	p.State = PostStateModerate
	return p
}

// SetStatePublic - setter for state to write
func (p *PostInBlog) SetStatePublic() *PostInBlog {
	p.State = PostStatePublic
	return p
}

// SetStateBlocked - setter for state to write
func (p *PostInBlog) SetStateBlocked() *PostInBlog {
	p.State = PostStateBlocked
	return p
}

// SetCreatedAt - setter for ID
func (p *PostInBlog) SetCreatedAt(createdAt string) *PostInBlog {
	// TODO: add time.Parse and check, set to Now or something else. You must think about it!
	p.CreatedAt = createdAt
	return p
}

// SetModifiedAt - setter for ModifiedAt
func (p *PostInBlog) SetModifiedAt(modifiedAt string) *PostInBlog {
	// TODO: add time.Parse and check, set to Now or something else. You must think about it!
	p.ModifiedAt = modifiedAt
	return p

}

// SetParentPostID - setter for ParentPost
func (p *PostInBlog) SetParentPostID(pid string) *PostInBlog {
	p.ParentPostID = pid
	return p
}

// IncCountOfViews - increment of CountOfViews
func (p *PostInBlog) IncCountOfViews() *PostInBlog {
	p.CountOfViews++
	return p
}

// IncCountOfStars - increment of CountOfViews
func (p *PostInBlog) IncCountOfStars() *PostInBlog {
	p.CountOfStars++
	return p
}

// DecCountOfStars - decrement of CountOfViews
func (p *PostInBlog) DecCountOfStars() *PostInBlog {
	p.CountOfStars--
	return p
}

// SetTags - setter for ID
func (p *PostInBlog) SetTags(tags Tags) *PostInBlog {
	p.Tags = tags
	return p
}

// AddCommentsID - setter for ID
func (p *PostInBlog) AddCommentsID(commentIDs []string) *PostInBlog {
	// TODO: add check for exists commentID in the slice
	for _, commentID := range commentIDs {
		p.CommentsIDs = append(p.CommentsIDs, commentID)
	}
	return p
}

// [BusinessRules for Posts]

// GetTemplatePost - returns empty template with filled specified properties
func (p *PostInBlog) GetTemplatePost() PostInBlog {
	var template PostInBlog
	template.State = PostStatePublic // for testing set default public -> //PostStateModerate // all new post must be in moderate state
	template.Author.ID = AnonimousID // write now allow anonimous write posts
	return template
}

// PostAdd - adding new PostInBlog, returns new PostInBlog ID and error
// func (p *PostInBlog) PostAdd(newpost PostInBlog) error {
// 	// all new posts must be moderated
// 	p.State = PostStateModerate

// }

// User is any one who visit my blog
type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Nick       string `json:"nick"`
	EMail      string `json:"email"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
	UserRole   int    `json:"userrole"`
	Salt       string `json:"salt"`
	Avatar     string `json:"avatar"`
	// and more other properties
}

// UserRepository is a storage of Users
type UserRepository interface {
	Store(u User) (string, error)
	FindByToken(t string) (User, error)
	FindByID(id string) (User, error)
	Find() ([]User, error)
	Update(u User) error
	Delete(u User) error
}

// isAdmin - checks admin privileges
func (ur *User) isAdmin() bool {
	return ur.UserRole == UserAdmin
}

// Rubric is topic or headline of Post
type Rubric struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Tags - slice of labels/Tags
type Tags []string

// CommentOfPost is single comment for some Post
type CommentOfPost struct {
	ID           string `json:"id"`
	Author       User   `json:"author"`
	Content      string `json:"content"`
	CountOfStars int64  `json:"count_of_stars"`
	PostID       string `json:"postid"`
}

// CommentsOfPost is slice of comments
type CommentsOfPost []CommentOfPost

// CommentsRepository is a storage for comments maybe another storage and full text search
type CommentsRepository interface {
	Store(c CommentOfPost) (string, error)
	FindByID(is string) (CommentOfPost, error)
	FindByPostID(pid string) ([]CommentsOfPost, error)
	Update(c CommentOfPost) error
	Delete(c CommentOfPost) error
}
