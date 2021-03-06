// TODO: use google's context.Context

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pborman/uuid"

	"github.com/julienschmidt/httprouter"
)

// Req is a HTTP request wrapper giving you easy access to the params
// content type, and much more.
type Req struct {
	ID          string // The request id
	Response    http.ResponseWriter
	Request     *http.Request
	Params      *Params // Parameters from URL and form (including multipart). Keep in ctx instead?
	ContentType string  // Content-Type of the request
	body        []byte
}

func NewReq(w http.ResponseWriter, r *http.Request, p *Params) *Req {
	return &Req{
		ID:       uuid.New(),
		Response: &statusResponseWriter{w, 0},
		Request:  r,
		Params:   p,
	}
}

// WrapReq wraps a standard request.
func WrapReq(w http.ResponseWriter, r *http.Request) *Req {
	return NewReq(w, r, new(Params))
}

// WrapReq wraps a request dispatched by httprouter.
func WrapHttpRouterReq(w http.ResponseWriter, r *http.Request, ps httprouter.Params) *Req {
	params := new(Params)

	if ps != nil && len(ps) > 0 {
		params.Path = make(Values, len(ps))
		for _, p := range ps {
			params.Path[":"+p.Key] = splitValues([]string{p.Value}, ",")
		}
	}

	return NewReq(w, r, params)
}

// ResponseStatus returns the response status code if available yet (0 otherwise).
func (r *Req) ResponseStatus() int {
	var status int
	if srw, ok := r.Response.(*statusResponseWriter); ok {
		status = srw.status
	}
	return status
}

// ResolveContentType extracts content type from the request.
func (r *Req) ResolveContentType() string {
	contentType := r.Request.Header.Get("Content-Type")
	if contentType == "" {
		r.ContentType = "text/html"
	} else {
		r.ContentType = strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	}

	return r.ContentType
}

// JsonBody() extracts the body from a request as a byte array
// so it cam be unmarshalled if desired.
func (r *Req) JsonBody() ([]byte, error) {
	if r.Request.Body == nil {
		return []byte{}, nil
	}
	defer r.Request.Body.Close()
	return ioutil.ReadAll(r.Request.Body)
}

func (r *Req) JsonForm() ([]byte, error) {
	if r.Params.Form == nil {
		return []byte{}, nil
	}
	return json.Marshal(r.Params.Form)
}

// Decode decodes a request body into the value pointed to by v.
func (r *Req) Decode(v interface{}) error {
	if r.body == nil {
		b, err := ioutil.ReadAll(r.Request.Body)
		if err != nil {
			return err
		}
		r.body = b
	}

	return json.Unmarshal(r.body, v)
}

// handlePanic is a function usually used in defer
// to catch panics and return a 500 instead so the web server
// doesn't crash.
func (r *Req) handlePanic() {
	var err *Error
	if rec := recover(); rec != nil {
		switch val := rec.(type) {
		case error:
			err = WrapErr(val, http.StatusInternalServerError)
		default:
			err = WrapErr(fmt.Errorf("%v", val), http.StatusInternalServerError)
		}

		http.Error(r.Response, err.HTTPBody(), err.HTTPStatus())
	}
}

// NoContent sends a response with no body and a status code.
func (r *Req) NoContent(code int) {
	r.Response.WriteHeader(code)
}

// Redirect redirects the request using http.Redirect with status code.
func (r *Req) Redirect(code int, url string) {
	http.Redirect(r.Response, r.Request, url, code)
}
