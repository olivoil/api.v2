package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/antonholmquist/jason"
	"github.com/bmizerany/pat"
	"github.com/julienschmidt/httprouter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Login middleware
func login(r *Req) (error, int) {
	r.Set("userID", "1")
	return nil, http.StatusOK
}

// Auth middleware
func auth(r *Req) (error, int) {
	if _, ok := r.GetOk("userID"); ok {
		return nil, http.StatusOK
	}
	return errors.New("not authenticated"), http.StatusUnauthorized
}

type pet struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func makeAPI() *API {
	api := New("v1")

	api.Add(Endpoint{
		Verb:       "POST",
		Path:       "/pets",
		Middleware: []Middleware{MiddlewareFunc(auth)},
		Implementation: func(r *Req) {
			pet := pet{}
			err := r.Decode(&pet)

			b, err := json.Marshal(pet)
			if err != nil {
				panic(err)
			}

			r.Response.Header().Set("Content-Type", "application/json")
			r.Response.WriteHeader(200)
			r.Response.Write(b)
		},
	})

	api.Add(Endpoint{
		Verb:       "PUT",
		Path:       "/pets/:id",
		Middleware: []Middleware{MiddlewareFunc(login), MiddlewareFunc(auth)},
		Implementation: func(r *Req) {
			pet := pet{ID: r.Params.Get(":id")}
			err := r.Decode(&pet)

			b, err := json.Marshal(pet)
			if err != nil {
				panic(err)
			}

			r.Response.Header().Set("Content-Type", "application/json")
			r.Response.WriteHeader(200)
			r.Response.Write(b)
		},
	})

	return api
}

var _ = Describe("api integration", func() {
	Context("httprouter", func() {
		It("handles a failing middleware", func() {
			router := httprouter.New()
			api := makeAPI()
			api.Activate(router)

			res := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/v1/pets", strings.NewReader(`{"name":"moufassa","id":"abc"}`))
			router.ServeHTTP(res, req)

			Expect(res.Code).To(Equal(401))

			data, err := jason.NewObjectFromReader(res.Body)
			Expect(err).ToNot(HaveOccurred())

			slice, err := data.GetObjectArray("errors")
			Expect(err).ToNot(HaveOccurred())
			Expect(slice).ToNot(BeNil())

			msg := slice[0]
			Expect(msg).ToNot(BeNil())

			i, err := msg.GetString("status")
			Expect(err).ToNot(HaveOccurred())
			Expect(i).To(Equal("401"))

			v, err := msg.GetString("title")
			Expect(err).ToNot(HaveOccurred())
			Expect(v).To(Equal("not authenticated"))
		})

		It("works", func() {
			router := httprouter.New()
			api := makeAPI()
			api.Activate(router)

			res := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/v1/pets/abc", strings.NewReader(`{"name":"simba"}`))
			router.ServeHTTP(res, req)

			Expect(res.Code).To(Equal(200))

			data, err := jason.NewObjectFromReader(res.Body)
			Expect(err).ToNot(HaveOccurred())

			name, err := data.GetString("name")
			Expect(err).ToNot(HaveOccurred())
			Expect(name).To(Equal("simba"))

			id, err := data.GetString("id")
			Expect(err).ToNot(HaveOccurred())
			Expect(id).To(Equal("abc"))
		})
	})

	Context("Gorilla Pat", func() {
		It("handles a failing middleware", func() {
			router := pat.New()
			api := makeAPI()
			api.Activate(router)

			res := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/v1/pets", strings.NewReader(`{"name":"moufassa","id":"abc"}`))
			router.ServeHTTP(res, req)

			Expect(res.Code).To(Equal(401))

			data, err := jason.NewObjectFromReader(res.Body)
			Expect(err).ToNot(HaveOccurred())

			slice, err := data.GetObjectArray("errors")
			Expect(err).ToNot(HaveOccurred())
			Expect(slice).ToNot(BeNil())

			msg := slice[0]
			Expect(msg).ToNot(BeNil())

			i, err := msg.GetString("status")
			Expect(err).ToNot(HaveOccurred())
			Expect(i).To(Equal("401"))

			v, err := msg.GetString("title")
			Expect(err).ToNot(HaveOccurred())
			Expect(v).To(Equal("not authenticated"))
		})

		It("works", func() {
			router := pat.New()
			api := makeAPI()
			api.Activate(router)

			res := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/v1/pets/abc", strings.NewReader(`{"name":"simba"}`))
			router.ServeHTTP(res, req)

			Expect(res.Code).To(Equal(200))

			data, err := jason.NewObjectFromReader(res.Body)
			Expect(err).ToNot(HaveOccurred())

			name, err := data.GetString("name")
			Expect(err).ToNot(HaveOccurred())
			Expect(name).To(Equal("simba"))

			id, err := data.GetString("id")
			Expect(err).ToNot(HaveOccurred())
			Expect(id).To(Equal("abc"))
		})
	})
})
