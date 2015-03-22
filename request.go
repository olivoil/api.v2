package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"code.google.com/p/go-uuid/uuid"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

// Req is a HTTP request wrapper giving you easy access to the params
// content type, and much more.
type Req struct {
	ID          string // The request id
	Response    http.ResponseWriter
	Request     *http.Request
	Params      *Params // Parameters from URL and form (including multipart).
	ContentType string

	m   sync.Mutex
	Log []string // A slice of log messages attached to the request.

	decoder *json.Decoder
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
	if len(ps) > 0 {
		vars := map[string]string{}
		for _, p := range ps {
			vars[p.Key] = p.Value
		}
		registerVars(r, vars)
	}

	return NewReq(w, r, new(Params))
}

// registerVars adds the matched route variables to the URL query.
func registerVars(r *http.Request, vars map[string]string) {
	parts, i := make([]string, len(vars)), 0
	for key, value := range vars {
		parts[i] = url.QueryEscape(":"+key) + "=" + url.QueryEscape(value)
		i++
	}
	q := strings.Join(parts, "&")
	if r.URL.RawQuery == "" {
		r.URL.RawQuery = q
	} else {
		r.URL.RawQuery += "&" + q
	}
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
	json, err := ioutil.ReadAll(r.Request.Body)
	if err != nil {
		return json, err
	}
	return json, nil
}

// Decode decodes a request body into the value pointed to by v.
func (r *Req) Decode(v interface{}) error {
	if r.decoder == nil {
		r.decoder = json.NewDecoder(r.Request.Body)
	}

	return r.decoder.Decode(v)
}

// handlePanic is a function usually used in defer
// to catch panics and return a 500 instead so the web server
// doesn't crash.
func (r *Req) handlePanic() {
	var err Error
	if rec := recover(); rec != nil {
		switch val := rec.(type) {
		case error:
			err = WrapErr(val, http.StatusInternalServerError)
		default:
			err = WrapErr(fmt.Errorf("%v", val), http.StatusInternalServerError)
		}

		http.Error(r.Response, err.HTTPBody(), err.HTTPStatus())
		r.AddLog(fmt.Sprintf("%s - %v\n", err.Error(), r.Request))
		// if err.Stack != nil {
		// 	r.AddLog(string(err.Stack))
		// }
	}
}

// Add a string to the request log
func (r *Req) AddLog(s string) {
	r.m.Lock()
	r.Log = append(r.Log, s)
	r.m.Unlock()
}

// Expose request context
func (r *Req) Set(k interface{}, v interface{}) {
	context.Set(r.Request, k, v)
}

func (r *Req) Get(k interface{}) interface{} {
	return context.Get(r.Request, k)
}

func (r *Req) GetOk(k interface{}) (interface{}, bool) {
	return context.GetOk(r.Request, k)
}

func (r *Req) Clear() {
	context.Clear(r.Request)
}
