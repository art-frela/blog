package infra

import (
	"strings"
	"testing"

	"github.com/art-frela/blog/domain"
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.NewEntry(&logrus.Logger{})
	repo   = &MongoPostRepo{
		mongoURL:       "mongodb://elk-01.watcom.local:27017",
		log:            logger,
		database:       "testblog",
		collectionName: "posts",
	}
)

func TestConnDB(t *testing.T) {
	tests := []struct {
		name    string
		salt    string
		wantErr bool
	}{
		{"succ-connect-case", "27017", false},
		{"wron-connect-case", "27018", true},
	}
	defer func() {
		repo.mongoURL = "mongodb://elk-01.watcom.local:27017"
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.mongoURL = strings.Replace(repo.mongoURL, "27017", tt.salt, -1)
			_, err := repo.connDB()
			if (err != nil) != tt.wantErr {
				t.Errorf("connDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFindByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"nodoc-findByID-case", "777", true},
	}
	session, err := repo.connDB()
	if err != nil {
		repo.log.Fatalf("connect to mongoDB error, %v", err)
	}
	repo.session = session
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.FindByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		name    string
		offset  int
		limit   int
		wantErr bool
	}{
		{"norm-find-case", 0, 10, false},
		{"zero-find-case", 0, 0, false},
	}
	session, err := repo.connDB()
	if err != nil {
		repo.log.Fatalf("connect to mongoDB error, %v", err)
	}
	repo.session = session
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.Find(tt.offset, tt.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSave(t *testing.T) {
	tests := []struct {
		name    string
		post    domain.PostInBlog
		wantErr bool
	}{
		{"first-insert", domain.PostInBlog{
			Title: "1st",
		}, false},
		{"second-insert", domain.PostInBlog{
			Title: "2nd",
		}, false},
	}
	session, err := repo.connDB()
	if err != nil {
		repo.log.Fatalf("connect to mongoDB error, %v", err)
	}
	repo.session = session
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.Save(tt.post)
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name    string
		doc     domain.PostInBlog
		wantErr bool
	}{
		{"first-update", domain.PostInBlog{
			ID:    "777",
			Title: "1stUPDATE",
		}, true},
	}
	session, err := repo.connDB()
	if err != nil {
		repo.log.Fatalf("connect to mongoDB error, %v", err)
	}
	repo.session = session
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(tt.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
