package api

type Middleware struct {
	// identifier
	Name string
	// documentation available from within the code for easier debugging
	Description string
	// Untyped params used by the middleware
	// each middleware is responsible to check its params and their types.
	Params interface{}
	// HandlerFunc to process the incoming request and
	// returns a http error code and error message if needed.
	Handler func(r *Req) (error, int)
}

type MiddlewareStack []Middleware
