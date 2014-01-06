package api

type Middleware interface {
	// identifier
	Name() string
	// documentation available from within the code for easier debugging
	Description() string
	// HandlerFunc to process the incoming request and
	// returns a http error code and error message if needed.
	Run(r *Req) (error, int)
}

type MiddlewareStack []Middleware
