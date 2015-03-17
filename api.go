package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Top-level API structure
type API struct {
	Endpoints []Endpoint
	options   map[string][]string
}

// HandlerFunc
type HandlerFunc func(*Req)

func (f HandlerFunc) Serve(req *Req) {
	f(req)
}

// Handler
type Handler interface {
	Serve(req *Req)
}

// Router
type Router interface {
	Add(method, path string, handler Handler)
}

func (api *API) Add(e Endpoint) {
	// add endpoint
	api.Endpoints = append(api.Endpoints, e)

	// collect options for each path
	if _, ok := api.options[e.Path]; !ok {
		api.options[e.Path] = []string{"OPTIONS"}
	}
	api.options[e.Path] = append(api.options[e.Path], e.Verb)
}

func (api *API) Activate(router Router) {
	for _, endpoint := range api.Endpoints {
		router.Add(endpoint.Verb, endpoint.Path, &endpoint)
	}

	for path, verbs := range api.options {
		router.Add("OPTIONS", path, HandlerFunc(func(r *Req) {
			r.Response.Header().Set("Allow", strings.Join(verbs, ","))
			r.Response.WriteHeader(http.StatusNoContent)
		}))
	}
}

func WrapRouter(v interface{}) (Router, error) {
	if r, ok := v.(*http.ServeMux); ok {
		return &httpServeMuxAdapter{r}, nil
	}

	if r, ok := v.(*httprouter.Router); ok {
		return &httprouterAdapter{r}, nil
	}

	return nil, errors.New(fmt.Sprintf("Cannot wrap unsupported router type: %+v", v))
}

type httpServeMuxAdapter struct {
	*http.ServeMux
}

func (router *httpServeMuxAdapter) Add(method, path string, h Handler) {
	router.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			req := WrapReq(w, r)
			h.Serve(req)
		}
	}))
}

type httprouterAdapter struct {
	*httprouter.Router
}

func (router *httprouterAdapter) Add(method, path string, h Handler) {
	router.Handle(method, path, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		req := WrapHttpRouterReq(w, r, ps)
		h.Serve(req)
	})
}
