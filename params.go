package api

import (
	"log"
	"mime/multipart"
	"net/url"
	"os"
)

// TODO : Empirically determine sane limit
const maxMultipartFormMemory = 10 << 20 // 10Mb

// Wrapper for the request params.
type Params struct {
	url.Values // A unified view of all the individual param maps below.

	Query url.Values // Parameters from the query string, e.g. /index?limit=10
	Form  url.Values // Parameters from the request body.

	Files          map[string][]*multipart.FileHeader // Files uploaded in a multipart form
	tmpFiles       []*os.File                         // Temp files used during the request.
	RequiredParams []string                           // The required that have to be present for the params to be valid
}

func (r *Req) ParseParams() error {
	r.Params.Query = r.Request.URL.Query()
	if len(r.ContentType) == 0 {
		r.ResolveContentType()
	}
	// Parse the body depending on the content type.
	switch r.ContentType {
	case "application/x-www-form-urlencoded":
		// Typical form.
		if err := r.Request.ParseForm(); err != nil {
			log.Printf("[HTTP 400] Error parsing request body: %s - %v\n", err.Error(), r.Request)
			return err
		} else {
			r.Params.Form = r.Request.Form
		}
	case "multipart/form-data":
		// Multipart form.
		if err := r.Request.ParseMultipartForm(maxMultipartFormMemory); err != nil {
			// We have a multipart error, do not delete tmp file that holds the body
			return err
		} else {
			r.Params.Form = r.Request.MultipartForm.Value
			r.Params.Files = r.Request.MultipartForm.File
		}
	}

	r.Params.Values = r.Params.calcValues()
	return nil
}

// calcValues returns a unified view of the component param maps.
func (p *Params) calcValues() url.Values {
	numParams := len(p.Query) + len(p.Form)

	// If there were no params, return an empty map.
	if numParams == 0 {
		return make(url.Values, 0)
	}

	// If only one of the param sources has anything, return that directly.
	switch numParams {
	case len(p.Query):
		return p.Query
	case len(p.Form):
		return p.Form
	}

	// Copy everything into the same map.
	values := make(url.Values, numParams)
	for k, v := range p.Query {
		values[k] = append(values[k], v...)
	}
	for k, v := range p.Form {
		values[k] = append(values[k], v...)
	}
	return values
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
