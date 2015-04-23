// TODO: use google's context.Context

package api

import (
	"fmt"

	"golang.org/x/net/context"
)

type RequestParser interface {
	ParseRequest(context.Context, *Req) (context.Context, error)
}

type ResponseMarshaller interface {
	Body(context.Context) interface{}
	Headers(context.Context) map[string]string
	Status(context.Context) int
}

type Resource struct {
	Req    *Req
	Source DataSource
}

// DataSource provides methods needed for CRUD.
type DataSource interface {
	// FindOne returns a model from a parsed query
	FindOne(context.Context) (context.Context, error)

	// FindAll returns all objects specified in query
	FindAll(context.Context) (context.Context, error)

	// Create a new object and return its ID
	Create(context.Context) (context.Context, error)

	// Update an object and return its ID
	Update(context.Context) (context.Context, error)

	// Delete an object
	Delete(context.Context) (context.Context, error)
}

func NewResource(req *Req, source DataSource) *Resource {
	return &Resource{
		Source: source,
		Req:    req,
	}
}

func (r *Resource) HandleIndex(ctx context.Context, rp RequestParser) (ct context.Context, err error) {
	c, err := rp.ParseRequest(ctx, r.Req)
	if err != nil {
		return
	}

	ct, err = r.Source.FindAll(c)
	return
}

func (r *Resource) HandleRead(ctx context.Context, rp RequestParser) (ct context.Context, err error) {
	c, err := rp.ParseRequest(ctx, r.Req)
	if err != nil {
		return
	}

	ct, err = r.Source.FindOne(c)
	return
}

func (r *Resource) HandleCreate(ctx context.Context, rp RequestParser) (ct context.Context, err error) {
	// Unmarshal request model into model values
	c, err := rp.ParseRequest(ctx, r.Req)
	if err != nil {
		return
	}

	cx, err := r.Source.Create(c)
	if err != nil {
		return
	}

	ct, err = r.Source.FindOne(cx)
	return
}

func (r *Resource) HandleUpdate(ctx context.Context, rp RequestParser) (ct context.Context, err error) {
	c, err := rp.ParseRequest(ctx, r.Req)
	if err != nil {
		return
	}

	cx, err := r.Source.Update(c)
	if err != nil {
		return
	}

	ct, err = r.Source.FindOne(cx)
	return
}

func (r *Resource) HandleDelete(ctx context.Context, rp RequestParser) (ct context.Context, err error) {
	c, err := rp.ParseRequest(ctx, r.Req)
	if err != nil {
		return
	}

	cx, err := r.Source.FindOne(c)
	if err != nil {
		return
	}

	ct, err = r.Source.Delete(cx)
	return
}

func (r *Resource) HandleError(err error) {
	handleError(r.Req, err)
}

func HandleError(req *Req, err error) {
	handleError(req, err)
}

func handleError(req *Req, err error) {
	apiErr, ok := err.(*Error)
	if !ok {
		apiErr = WrapErr(err, 500)
	}

	req.Response.WriteHeader(apiErr.HTTPStatus())
	fmt.Fprintln(req.Response, apiErr.HTTPBody())
}

func (r *Resource) Send(ctx context.Context, rm ResponseMarshaller) error {
	return send(ctx, r.Req, rm)
}

func Send(ctx context.Context, req *Req, rm ResponseMarshaller) error {
	return send(ctx, req, rm)
}

func send(ctx context.Context, req *Req, rm ResponseMarshaller) error {
	encoder := JsonEncoder{}
	data, err := encoder.Encode(rm.Body(ctx))
	if err != nil {
		return err
	}
	headers := rm.Headers(ctx)
	if headers != nil {
		for key, val := range headers {
			req.Response.Header().Set(key, val)
		}
	}
	req.Response.WriteHeader(rm.Status(ctx))
	req.Response.Write(data)
	return nil
}
