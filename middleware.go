package api

import (
	"reflect"
	"runtime"

	"golang.org/x/net/context"
)

type Middleware interface {

	// HandlerFunc to process the incoming request and
	// returns a http error code and error message if needed.
	Run(ctx context.Context, r *Req) (context.Context, error)

	// Name of the middleware for debugging
	Name() string
}

// MiddlewareFunc transforms a function with the right signature
// into a Middleware
type MiddlewareFunc func(ctx context.Context, r *Req) (context.Context, error)

func (m MiddlewareFunc) Run(ctx context.Context, r *Req) (context.Context, error) {
	return m(ctx, r)
}

func (m MiddlewareFunc) Name() string {
	return runtime.FuncForPC(reflect.ValueOf(m).Pointer()).Name()
}

type MiddlewareStack []Middleware
