// The api package is a Go package designed to provide tools to handle API requests.
// The primarily goal is to support HTTP API endpoints.
package api

// Top-level API structure
type API struct {
	Endpoints     []Endpoint
	Documentation *APIDoc
}

func (api *API) Add(e Endpoint) {
	api.Endpoints = append(api.Endpoints, e)
}

type APIDoc struct {
	Swagger  string
	Info     *Info
	Host     string
	BasePath string
	Schemes  []string
	Consumes []string
	Produces []string
}

func New(options Options) API {
	return API{}
}

type Options struct {
	Host           string
	BasePath       string
	Schemes        []string
	Consumes       []string
	Produces       []string
	Version        string
	Title          string
	Description    string
	TermsOfService string
	Contact        *Contact
	License        *License
}

func (api API) MarshalJSON() ([]byte, error) {
	return nil, nil
}
