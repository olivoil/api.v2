package api

import "net/http"

// An Endpoint is the structure representing an API endpoint
// encapsulating information needed to
// dispatch a request and generate documentation.
type Endpoint struct {
	Url  string
	Verb string

	// The middlewares to execute on the request.
	Middleware MiddlewareStack

	// Called after middleware stack was executed on the request
	Implementation func(r *Req)

	// Used to generate API documentation
	Documentation *Operation
}

// Append a middleware to the middleware stack.
func (e *Endpoint) Use(mw Middleware) {
	e.Middleware = append(e.Middleware, mw)
}

// ServeHTTP implements the http.Handler interface
func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create the context
	req := WrapReq(w, r)
	defer req.handlePanic()

	// Parse the parameters and cleanup
	defer cleanUpParams(req)
	err := req.ParseParams()

	// We must return a 400 and stop here if there was a problem parsing the request.
	if err != nil {
		e := WrapErr(err, 400)
		http.Error(req.Response, e.HTTPBody(), e.HTTPStatus())
		return
	}

	// call each middleware
	for _, m := range e.Middleware {
		err, code := m.Run(req)
		if err != nil {
			e := WrapErr(err, code)
			http.Error(req.Response, e.HTTPBody(), e.HTTPStatus())
		}
	}

	// Dispatch the request via the endpoint
	e.Implementation(req)
}

// HandlerFunc converts an Endpoint to a http.HandlerFunc
func (e *Endpoint) HandlerFunc() http.HandlerFunc {
	return http.HandlerFunc(e.ServeHTTP)
}
