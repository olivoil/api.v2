// TODO: extract to own package (i.e. go-swagger)
package api

import (
	"encoding/json"
	"errors"
	"strconv"
)

type Documentation struct {
	Swagger  string   `json:"swagger,required"`
	Info     *Info    `json:"info,required"`
	Host     string   `json:"host,omitempty"`
	BasePath string   `json:"basePath,omitempty"`
	Schemes  []string `json:"schemes,omitempty" enum:"http,https,ws,wss"`
	Consumes []string `json:"consumes,omitempty"`
	Produces []string `json:"produces,omitempty"`
	Paths    []*Path  `json:"paths,required"`
	Tags     []*Tag   `json:"tags,omitemtpy"`
}

func (d *Documentation) AddEndpoint(e Endpoint) error {
	pathItem := PathItem{}

	switch e.Verb {
	case "GET":
		pathItem.Get = e.Documentation
	case "POST":
		pathItem.Post = e.Documentation
	case "PUT":
		pathItem.Put = e.Documentation
	case "DELETE":
		pathItem.Delete = e.Documentation
	case "PATCH":
		pathItem.Patch = e.Documentation
	case "OPTION":
		pathItem.Options = e.Documentation
	case "HEAD":
		pathItem.Head = e.Documentation
	default:
		return errors.New("VERB not recognized: " + e.Verb)
	}

	path := Path{
		RelativeUrl: e.Url,
		PathItem:    &pathItem,
	}

	d.Paths = append(d.Paths, &path)
	return nil
}

type Path struct {
	RelativeUrl string
	PathItem    *PathItem
}

type Tag struct {
	Name                  string                 `json:"name,required"`
	Description           string                 `json:"description,omitempty"`
	ExternalDocumentation *ExternalDocumentation `json:"externalDocs,omitempty"`
}

type ExternalDocumentation struct {
	Url         string `json:"url,required"`
	Description string `json:"description,omitempty"`
}

func (p Path) MarshalJSON() ([]byte, error) {
	var obj map[string]*PathItem
	obj[p.RelativeUrl] = p.PathItem
	return json.Marshal(obj)
}

type PathItem struct {
	Get        *Operation   `json:"get,omitempty"`
	Post       *Operation   `json:"post,omitempty"`
	Put        *Operation   `json:"put,omitempty"`
	Delete     *Operation   `json:"delete,omitempty"`
	Options    *Operation   `json:"options,omitempty"`
	Head       *Operation   `json:"head,omitempty"`
	Patch      *Operation   `json:"patch,omitempty"`
	Parameters []*Parameter `json:"parameters,omitempty"`
}

type Operation struct {
	ID          string       `json:"operationId,omitempty"`
	Tags        []string     `json:"tags,omitempty"`
	Summary     string       `json:"summary,omitempty"`
	Description string       `json:"description,omitempty"`
	Schemes     []string     `json:"schemes,omitempty"`
	Consumes    []string     `json:"consumes,omitempty"`
	Produces    []string     `json:"produces,omitempty"`
	Parameters  []*Parameter `json:"parameters,omitempty"`
	Responses   *Responses   `json:"responses,omitempty"`
	Deprecated  bool         `json:"deprecated,omitempty"`
	Security    []*Security  `json:"security,omitempty"`
}

type Parameter struct {
	Name        string `json:"name,required"`
	In          string `json:"in,required" enum:"query,header,path,formData,body"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`

	Schema *Schema `json:"schema,omitempty"`

	Type             string `json:"type,omitempty" enum:"string,number,integer,boolean,array,file"`
	Format           string `json:"format,omitempty"`
	Items            string `json:"items,omitempty"`
	CollectionFormat string `json:"collectionFormat,omitempty" enum:"csv,ssv,tsv,pipes,multi"`
}

type Responses []*Response

func (r Responses) MarshalJSON() ([]byte, error) {
	var rs map[string]*ResponseItem

	for _, resp := range r {
		rs[strconv.Itoa(resp.Status)] = resp.ResponseItem
	}

	return json.Marshal(rs)
}

type Response struct {
	Status       int
	ResponseItem *ResponseItem
}

type ResponseItem struct {
	Description string   `json:"description,required"`
	Schema      *Schema  `json:"schema,omitempty"`
	Headers     *Headers `json:"headers,omitempty"`
}

type Security struct {
}

type Schema struct {
}

type Item struct {
}

type Headers []*Header

func (h Headers) MarshalJSON() ([]byte, error) {
	var hs map[string]*HeaderItem

	for _, header := range h {
		hs[header.Name] = header.HeaderItem
	}

	return json.Marshal(hs)
}

type Header struct {
	Name       string
	HeaderItem *HeaderItem
}

type HeaderItem struct {
	Type             string  `json:"type,required" enum:string,number,integer,boolean,array"`
	Format           string  `json:"format,omitempty"`
	Description      string  `json:"description,omitempty"`
	CollectionFormat string  `json:"collectionFormat,omitempty" enum:"csv,ssv,tsv,pipes,multi"`
	Items            []*Item `json:"items,omitempty"`
}

type Info struct {
	Version        string   `json:"version,required"`
	Title          string   `json:"title,required"`
	Description    string   `json:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
}

type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	Name string `json:"name,required"`
	URL  string `json:"url,omitempty"`
}
