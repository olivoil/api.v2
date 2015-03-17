package api

type Middleware interface {

	// HandlerFunc to process the incoming request and
	// returns a http error code and error message if needed.
	Run(r *Req) (error, int)
}

type MiddlewareStack []Middleware
