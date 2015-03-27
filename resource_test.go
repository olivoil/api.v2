package api

import (
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/julienschmidt/httprouter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/guregu/null.v2"
)

type Post struct {
	ID          string
	Title       string
	Value       null.Float
	Author      *User
	AuthorID    string
	Comments    []Comment
	CommentsIDs []string
}

type Comment struct {
	ID    string
	Value string
}

type User struct {
	ID   string
	Name string
}

type fixtureSource struct {
	posts map[string]*Post
}

func (s *fixtureSource) FindAll(req *Req) (interface{}, error) {
	var postsSlice []Post

	if limit := req.Params.Get("limit"); limit != "" {
		if l, err := strconv.ParseInt(limit, 10, 64); err == nil {
			postsSlice = make([]Post, l)
			length := len(s.posts)
			for i := 0; i < length; i++ {
				postsSlice[i] = *s.posts[strconv.Itoa(i+1)]
				if i+1 >= int(l) {
					break
				}
			}
		} else {
			return nil, err
		}
	} else {
		postsSlice = make([]Post, len(s.posts))
		length := len(s.posts)
		for i := 0; i < length; i++ {
			postsSlice[i] = *s.posts[strconv.Itoa(i+1)]
		}
	}

	return postsSlice, nil
}

func (s *fixtureSource) FindOne(id string, req *Req) (interface{}, error) {
	if p, ok := s.posts[id]; ok {
		return *p, nil
	}
	return nil, NewError(http.StatusNotFound, "post not found")
}

func (s *fixtureSource) FindMultiple(IDs []string, req *Req) (interface{}, error) {
	var posts []Post

	for _, id := range IDs {
		if p, ok := s.posts[id]; ok {
			posts = append(posts, *p)
		}
	}

	if len(posts) > 0 {
		return posts, nil
	}

	return nil, NewError(http.StatusNotFound, "post not found")
}

func (s *fixtureSource) Create(obj interface{}, req *Req) (string, error) {
	p := obj.(Post)

	if p.Title == "" {
		err := NewError(http.StatusBadRequest, "Missing title")
		return "", err.Add(Error{ID: "SomeErrorID", Path: "Title"})
	}

	maxID := 0
	for k := range s.posts {
		id, _ := strconv.Atoi(k)
		if id > maxID {
			maxID = id
		}
	}
	newID := strconv.Itoa(maxID + 1)
	p.ID = newID
	s.posts[newID] = &p
	return newID, nil
}

func (s *fixtureSource) Delete(id string, req *Req) error {
	delete(s.posts, id)
	return nil
}

func (s *fixtureSource) Update(obj interface{}, req *Req) error {
	p := obj.(Post)
	if oldP, ok := s.posts[p.ID]; ok {
		oldP.Title = p.Title
		return nil
	}
	return NewError(http.StatusNotFound, "post not found")
}

func makeSource() *fixtureSource {
	return &fixtureSource{map[string]*Post{
		"1": &Post{
			ID:    "1",
			Title: "Hello, World!",
			Author: &User{
				ID:   "1",
				Name: "Dieter",
			},
			Comments: []Comment{Comment{
				ID:    "1",
				Value: "This is a stupid post!",
			}},
		},
		"2": &Post{ID: "2", Title: "I am NR. 2"},
		"3": &Post{ID: "3", Title: "I am NR. 3"},
	}}
}

var _ = Describe("Resource", func() {

	It("HandleIndex", func() {
		source := makeSource()
		res := NewResource(Post{}, source)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/posts?limit=2", nil)
		req := WrapReq(w, r)
		err := req.ParseParams()
		Expect(err).ToNot(HaveOccurred())

		result, err := res.HandleIndex(req)
		posts := result.([]Post)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(posts)).To(Equal(2))
	})

	It("HandleRead", func() {
		source := makeSource()
		res := NewResource(Post{}, source)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/posts/2", nil)
		req := WrapHttpRouterReq(w, r, httprouter.Params{httprouter.Param{Key: "id", Value: "2"}})
		err := req.ParseParams()
		Expect(err).ToNot(HaveOccurred())

		result, err := res.HandleRead(req)
		Expect(err).ToNot(HaveOccurred())
		post, ok := result.(Post)
		Expect(ok).To(Equal(true))
		Expect(post.ID).To(Equal("2"))
	})

	PIt("HandleCreate", func() {
	})

	PIt("HandleUpdate", func() {
	})

	PIt("HandleDelete", func() {
	})

	PIt("HandleError", func() {
	})
})
