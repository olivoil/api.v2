package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Req is a HTTP request wrapper giving you easy access to the params
// content type, and much more.
type Req struct {
	Request     *http.Request
	Response    http.ResponseWriter
	Params      *Params // Parameters from URL and form (including multipart).
	ContentType string
	UserId      int
	Log         []string // A slice of log messages attached to the request.
	Id          string   // The request id
}

// NewReq() automatically wraps a standard request.
func WrapReq(w http.ResponseWriter, r *http.Request) *Req {
	return &Req{
		Request:  r,
		Response: &statusResponseWriter{w, 0},
		Params:   new(Params),
	}
}

// ResponseStatus() returns the response status code if available yet (0 otherwise).
func (r *Req) ResponseStatus() int {
	var status int
	if srw, ok := r.Response.(*statusResponseWriter); ok {
		status = srw.status
	}
	return status
}

// ResolveContentType() extracts content type from the request.
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
	json, err := ioutil.ReadAll(r.Request.Body)
	if err != nil {
		return json, err
	}
	return json, nil
}

// handlePanic() is a function usually used in defer
// to catch panics and return a 500 instead so the web server
// doesn't crash.
func (r *Req) handlePanic() {
	var err error
	if rec := recover(); rec != nil {
		r.Response.WriteHeader(http.StatusInternalServerError)
		// TODO: customize the error response based on the content type.
		r.Response.Write([]byte("Something went wrong :("))
		// The recovered panic may not be an error
		switch val := rec.(type) {
		case error:
			err = val
		default:
			err = fmt.Errorf("%v", val)
		}
		r.Log = append(r.Log, fmt.Sprintf("[HTTP 500] %s - %v\n", err, r.Request))
	}
}
