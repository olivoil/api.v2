package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

func New(prefix string) *API {
	return &API{
		Prefix:    prefix,
		Endpoints: []Endpoint{},
		options:   map[string][]string{},
	}
}

// Top-level API structure
type API struct {
	Endpoints  []Endpoint
	Middleware MiddlewareStack
	options    map[string][]string
	Prefix     string
}

// HandlerFunc
type HandlerFunc func(context.Context, *Req)

func (f HandlerFunc) Serve(ctx context.Context, req *Req) {
	f(ctx, req)
}

// Handler
type Handler interface {
	Serve(ctx context.Context, req *Req)
}

// Router is an interface that helps activating an API
// to different types of router libraries (pat, httprouter, http, etc.)
type Router interface {
	Add(method, path string, handler Handler)
}

// Add() adds an endpoint to the API
func (api *API) Add(e Endpoint) {
	// add endpoint
	api.Endpoints = append(api.Endpoints, e)

	if api.options == nil {
		api.options = map[string][]string{}
	}

	// collect options for each path
	if _, ok := api.options[e.Path]; !ok {
		api.options[e.Path] = []string{"OPTIONS"}
	}
	api.options[e.Path] = append(api.options[e.Path], e.Method)
}

// Append a middleware to the middleware stack.
func (api *API) Use(mw Middleware) {
	api.Middleware = append(api.Middleware, mw)
}

// Activate() registers all endpoints in the api
// to the provided router
func (api *API) Activate(r interface{}) error {
	router, err := WrapRouter(r)
	if err != nil {
		return err
	}

	for _, endpoint := range api.Endpoints {
		api.activateEndpoint(endpoint, router)
	}

	for path, verbs := range api.options {
		router.Add("OPTIONS", api.Prefix+path, HandlerFunc(func(ctx context.Context, r *Req) {
			r.Response.Header().Set("Allow", strings.Join(verbs, ","))
			r.Response.WriteHeader(http.StatusNoContent)
		}))
	}

	return nil
}

func (api *API) activateEndpoint(e Endpoint, r Router) {
	r.Add(e.Method, api.Prefix+e.Path, HandlerFunc(func(ctx context.Context, r *Req) {
		for _, m := range api.Middleware {
			c, err := m.Run(ctx, r)
			if err != nil {
				er := WrapErr(err, 0)
				r.Response.WriteHeader(er.HTTPStatus())
				fmt.Fprintln(r.Response, er.HTTPBody())
				return
			}
			ctx = c
		}

		e.Serve(ctx, r)
	}))
}

// Wrap a router to be used with Activate
// i.e. api.Activate(WrapRouter(router))
func WrapRouter(v interface{}) (Router, error) {
	if r, ok := v.(*httprouter.Router); ok {
		return &httprouterAdapter{r}, nil
	}

	if r, ok := v.(patRouter); ok {
		return &patAdapter{r: r}, nil
	}

	return nil, errors.New("Cannot wrap unsupported router type")
}

// Adapter for github.com/julienschmidt/httprouter
type httprouterAdapter struct {
	*httprouter.Router
}

func (router *httprouterAdapter) Add(method, path string, h Handler) {
	router.Handle(method, path, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		req := WrapHttpRouterReq(w, r, ps)
		h.Serve(context.Background(), req)
	})
}

// Adapter for pat-like routers (github.com/gorilla/pat, github.com/bmizerany/pat)
type patRouter interface {
	Add(meth, pat string, h http.Handler)
}

type patAdapter struct {
	r patRouter
}

func (router patAdapter) Add(method, path string, h Handler) {
	router.r.Add(method, path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := WrapReq(w, r)
		h.Serve(context.Background(), req)
	}))
}
