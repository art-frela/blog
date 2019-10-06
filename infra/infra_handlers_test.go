package infra

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPosts(t *testing.T) {
	tests := []struct {
		name     string
		offset   int
		limit    int
		httpCode int
	}{
		{"posts-zero-values", 0, 0, http.StatusOK},
		{"posts-empty-values", 50000, 1, http.StatusOK},
		{"posts-negativeOffset-values", -10, 10, http.StatusOK},
		{"posts-negativeLimit-values", 0, -10, http.StatusOK},
		{"posts-negative-values", -10, -5, http.StatusOK},
	}
	templatePATH = "../assets/templates/*.html"
	blogSRV := NewBlogServer(0, false)
	blogSRV.Run()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uri := fmt.Sprintf("/posts?offset=%d&limit=%d", tt.offset, tt.limit)
			req, err := http.NewRequest("GET", uri, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(blogSRV.controller.GetPosts)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.httpCode {
				t.Errorf("got http status: %d, expected %d", status, tt.httpCode)
			}
		})
	}

	blogSRV.Stop()
}

// func TestGetOnePost(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		id       string
// 		httpCode int
// 	}{
// 		{"post-exists-values", "5d90b1d3242abfd8fa7f8cc4", http.StatusOK},
// 		{"post-zero-values", "5d90b1d3242abfd8fa7f8cz4", http.StatusNotFound},
// 	}
// 	templatePATH = "../assets/templates/*.html"
// 	blogSRV := NewBlogServer("mongodb://elk-01.watcom.local:27017", 0, false)
// 	ws := blogSRV.Run(":8888")

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			uri := fmt.Sprintf("/posts/%s/", tt.id)
// 			req, err := http.NewRequest("GET", uri, nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			rr := httptest.NewRecorder()
// 			handler := http.HandlerFunc(blogSRV.controller.GetOnePost)
// 			handler.ServeHTTP(rr, req)

// 			if status := rr.Code; status != tt.httpCode {
// 				t.Errorf("got http status: %d, expected %d", status, tt.httpCode)
// 			}
// 		})
// 	}

// 	blogSRV.Stop(ws)
// }

func TestAddNewPost(t *testing.T) {
	tests := []struct {
		name     string
		post     string
		httpCode int
	}{
		{"post-new-post", `{"title":"1st"}`, http.StatusBadRequest},
		{"post-new2-post", `{"title":"2nd"}`, http.StatusBadRequest},
	}
	templatePATH = "../assets/templates/*.html"
	blogSRV := NewBlogServer(0, false)
	blogSRV.Run()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uri := fmt.Sprintf("/posts")
			reader := bytes.NewReader([]byte(tt.post))
			req, err := http.NewRequest("POST", uri, reader)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(blogSRV.controller.AddNewPost)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.httpCode {
				t.Errorf("got http status: %d, expected %d", status, tt.httpCode)
			}
		})
	}

	blogSRV.Stop()
}
