package api

import (
	"fmt"
	"net/http"
	"net/url"
)

type DocResp struct {
	StatusCode  int    // e.g. 200
	Body        string // example response body
	Description string // response description
}

// An Endpoint is the structure representing an API endpoint
// encapsulating information needed to dispatch a request and generate
// documentation.
type EndPoint struct {
	// The URL used to access the endpoint.
	Url *url.URL
	// GET, POST, PUT, etc.
	Method string
	// Description of the endpoint (goal, etc..)
	Description string
	// The middlewares to execute on the request.
	Middleware MiddlewareStack
	// documentation of a succesful response
	Successfull *DocResp
	// Documentation of a failed response
	Failed *DocResp
	//
	Implementation func(r *Req)
}

// Append a middleware to the middleware stack.
func (e *EndPoint) Use(mw *Middleware) {
	e.Middleware = append(e.Middleware, *mw)
}

//
func (e *EndPoint) Dispatch() http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Create the context
		req := WrapReq(w, r)
		defer req.handlePanic()
		// Parse the parameters and cleanup
		defer cleanUpParams(req)
		err := req.ParseParams()
		// We must return a 400 and stop here if there was a problem parsing the request.
		if err != nil {
			http.Error(req.Response, fmt.Sprintf("{\"error\":\"%s\"}", "Bad params"), 400)
			return
		}
		// call each middleware
		for _, m := range e.Middleware {
			err, code := m.Handler(req)
			if err != nil {
				http.Error(req.Response, fmt.Sprintf("{\"error\":\"%s\"}", err), code)
			}
		}
		// Dispatch the request via the endpoint
		e.Implementation(req)
	})

	return nil
}
