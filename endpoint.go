package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

// An Endpoint is the structure representing an API endpoint
// encapsulating information needed to
// dispatch a request and generate documentation.
type Endpoint struct {
	Verb string
	Path string

	// The middlewares to execute on the request.
	Middleware MiddlewareStack

	// Called after middleware stack was executed on the request
	Implementation func(ctx context.Context, r *Req)
}

// Append a middleware to the middleware stack.
func (e Endpoint) Use(mw Middleware) {
	e.Middleware = append(e.Middleware, mw)
}

// ServeHTTP implements the http.Handler interface
func (e Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := WrapReq(w, r)
	e.Serve(context.Background(), req)
}

// HandlerFunc converts an Endpoint to a http.HandlerFunc
func (e Endpoint) HandlerFunc() http.HandlerFunc {
	return http.HandlerFunc(e.ServeHTTP)
}

// Handle is a httprouter.Handle function
func (e Endpoint) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := WrapHttpRouterReq(w, r, ps)
	e.Serve(context.Background(), req)
}

// Serve dispatches an api.Req
func (e Endpoint) Serve(ctx context.Context, req *Req) {
	defer req.handlePanic()

	// Parse the parameters and cleanup
	defer cleanUpParams(req)
	err := req.ParseParams()

	// We must return a 400 and stop here if there was a problem parsing the request.
	if err != nil {
		er := WrapErr(err, 400)
		http.Error(req.Response, er.HTTPBody(), er.HTTPStatus())
		return
	}

	// call each middleware
	for _, m := range e.Middleware {
		err, c := m.Run(ctx, req)
		if err != nil {
			er := WrapErr(err, 0)
			req.Response.WriteHeader(er.HTTPStatus())
			fmt.Fprintln(req.Response, er.HTTPBody())
			return
		}
		ctx = c
	}

	// Dispatch the request via the endpoint
	e.Implementation(ctx, req)
}
