package api

import (
	"fmt"
	"mime/multipart"
	"os"
	"strings"
)

// TODO : Empirically determine sane limit
const maxMultipartFormMemory = 10 << 20 // 10Mb

// Wrapper for the request params.
type Params struct {
	Values // A unified view of all the individual param maps below.

	Query Values // Parameters from the query string, e.g. /index?limit=10
	Form  Values // Parameters from the request body.
	Path  Values // Parameters from the url path e.g. /users/:id

	Files          map[string][]*multipart.FileHeader // Files uploaded in a multipart form
	tmpFiles       []*os.File                         // Temp files used during the request.
	RequiredParams []string                           // The required that have to be present for the params to be valid
}

// Get implements the GetterSetter interface
func (p *Params) Get(s string) string {
	return p.Values.Get(s)
}

// GetAll implements the GetterSetter interface
func (p *Params) GetAll(s string) []string {
	return p.Values[s]
}

// Set implements the GetterSetter interface
func (p *Params) Set(s string, v string) {
	p.Values.Set(s, v)
}

// Append implements the GetterSetter interface
func (p *Params) Append(s string, v ...string) {
	if p.Values[s] == nil {
		p.Values[s] = v
		return
	}
	p.Values[s] = append(p.Values[s], v...)
}

// ParseParams from form, multipart, and query
func (r *Req) ParseParams() error {
	if len(r.ContentType) == 0 {
		r.ResolveContentType()
	}

	r.Params.Query = Values(r.Request.URL.Query())

	// Parse the body depending on the content type.
	switch r.ContentType {
	case "application/x-www-form-urlencoded":
		// Typical form.
		if err := r.Request.ParseForm(); err != nil {
			r.AddLog(fmt.Sprintf("Error parsing request body: %s", WrapErr(err, 400).Error()))
			return err
		} else {
			r.Params.Form = Values(r.Request.Form)
		}
	case "multipart/form-data":
		// Multipart form.
		if err := r.Request.ParseMultipartForm(maxMultipartFormMemory); err != nil {
			// We have a multipart error, do not delete tmp file that holds the body
			return err
		} else {
			r.Params.Form = Values(r.Request.MultipartForm.Value)
			r.Params.Files = r.Request.MultipartForm.File
		}
	}

	r.Params.Values = r.Params.calcValues()
	return nil
}

// calcValues returns a unified view of the component param maps.
func (p *Params) calcValues() Values {
	numParams := len(p.Query) + len(p.Form) + len(p.Path)

	// If there were no params, return an empty map.
	if numParams == 0 {
		return make(Values, 0)
	}

	// Copy everything into the same map.
	values := make(Values, numParams)
	for k, v := range p.Query {
		values.Append(k, splitValues(v, ",")...)
	}
	for k, v := range p.Form {
		values.Append(k, splitValues(v, ",")...)
	}
	for k, v := range p.Path {
		values.Append(k, splitValues(v, ",")...)
	}
	return values
}

func splitValues(vs []string, sep string) (res []string) {
	for _, v := range vs {
		res = append(res, strings.Split(v, sep)...)
	}
	return
}

func cleanUpParams(r *Req) error {
	// Delete temp files.
	if r.Request.MultipartForm != nil {
		err := r.Request.MultipartForm.RemoveAll()
		if err != nil {
			return err
		}
	}

	for _, tmpFile := range r.Params.tmpFiles {
		err := os.Remove(tmpFile.Name())
		if err != nil {
			return err
		}
	}

	return nil
}
