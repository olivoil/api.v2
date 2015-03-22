package api

import (
	"reflect"
	"runtime"
)

type Middleware interface {

	// HandlerFunc to process the incoming request and
	// returns a http error code and error message if needed.
	Run(r *Req) (error, int)

	// Name of the middleware for debugging
	Name() string
}

// MiddlewareFunc transforms a function with the right signature
// into a Middleware
type MiddlewareFunc func(r *Req) (error, int)

func (m MiddlewareFunc) Run(r *Req) (error, int) {
	return m(r)
}

func (m MiddlewareFunc) Name() string {
	return runtime.FuncForPC(reflect.ValueOf(m).Pointer()).Name()
}

type MiddlewareStack []Middleware
