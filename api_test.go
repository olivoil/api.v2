package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/antonholmquist/jason"
	"github.com/bmizerany/pat"
	"github.com/julienschmidt/httprouter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
)

type key int

var (
	userID  key = 0
	version key = 1
)

// API version middleware
func versionMiddleware(ctx context.Context, r *Req) (context.Context, error) {
	r.Response.Header().Set("X-Api-Version", "0.0.1")
	return context.WithValue(ctx, version, "0.0.1"), nil
}

// Login middleware
func login(ctx context.Context, r *Req) (context.Context, error) {
	c := context.WithValue(ctx, userID, "1")
	return c, nil
}

// Auth middleware
func auth(ctx context.Context, r *Req) (context.Context, error) {
	id, ok := ctx.Value(userID).(string)
	if !ok || id == "" {
		return ctx, NewError(http.StatusUnauthorized, "not authenticated")
	}
	return ctx, nil
}

type pet struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func makeAPI() *API {
	api := New("/v1")

	api.Use(MiddlewareFunc(versionMiddleware))

	api.Add(Endpoint{
		Method:     "PUT",
		Path:       "/pets/:id",
		Middleware: []Middleware{MiddlewareFunc(login), MiddlewareFunc(auth)},
		Implementation: func(ctx context.Context, r *Req) {
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

	api.Add(Endpoint{
		Method:     "POST",
		Path:       "/pets",
		Middleware: []Middleware{MiddlewareFunc(auth)},
		Implementation: func(ctx context.Context, r *Req) {
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

	return api
}

var _ = Describe("api integration", func() {
	Context("httprouter", func() {
		It("handles a failing middleware", func() {
			router := httprouter.New()
			api := makeAPI()
			api.Activate(router)

			res := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/v1/pets", strings.NewReader(`{"name":"moufassa","id":"king"}`))
			router.ServeHTTP(res, req)

			Expect(res.Header().Get("X-Api-Version")).To(Equal("0.0.1"))
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
			fmt.Printf("ServeHTTP %v %v\n", req.Method, req.URL.String())
			router.ServeHTTP(res, req)

			Expect(res.Header().Get("X-Api-Version")).To(Equal("0.0.1"))
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

			Expect(res.Header().Get("X-Api-Version")).To(Equal("0.0.1"))
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

			Expect(res.Header().Get("X-Api-Version")).To(Equal("0.0.1"))
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
