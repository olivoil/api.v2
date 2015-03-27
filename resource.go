package api

import (
	"errors"
	"fmt"
	"reflect"
)

type Resource struct {
	resourceType reflect.Type
	source       DataSource
}

// DataSource provides methods needed for CRUD.
type DataSource interface {
	// FindAll returns all objects
	FindAll(*Req) (interface{}, error)

	// FindOne returns an object by its ID
	FindOne(ID string, req *Req) (interface{}, error)

	// FindMultiple returns all objects for the specified IDs
	FindMultiple(IDs []string, req *Req) (interface{}, error)

	// Create a new object and return its ID
	Create(v interface{}, req *Req) (string, error)

	// Delete an object
	Delete(id string, req *Req) error

	// Update an object
	Update(obj interface{}, req *Req) error
}

// Request unmarshals a *Req into an interface value
type RequestModel interface {
	Unmarshal(req *Req) (interface{}, error)
}

func NewResource(model interface{}, source DataSource) *Resource {
	resourceType := reflect.TypeOf(model)
	if resourceType.Kind() != reflect.Struct {
		panic("pass an empty model struct to api.NewResource(Model{}, &DataSource{})!")
	}

	return &Resource{
		resourceType: resourceType,
		source:       source,
	}
}

func (r *Resource) HandleIndex(req *Req) (interface{}, error) {
	objs, err := r.source.FindAll(req)
	if err != nil {
		return nil, err
	}

	return objs, nil
}

func (r *Resource) HandleRead(req *Req) (interface{}, error) {
	ids := req.Params.GetAll(":id")

	var (
		obj interface{}
		err error
	)

	if len(ids) == 1 {
		obj, err = r.source.FindOne(ids[0], req)
	} else {
		obj, err = r.source.FindMultiple(ids, req)
	}

	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (r *Resource) HandleCreate(req *Req, mod RequestModel) (interface{}, error) {
	// Unmarshal request model into model values
	v, err := mod.Unmarshal(req)
	if err != nil {
		return nil, err
	}
	newObjs := reflect.ValueOf(v).Convert(reflect.SliceOf(r.resourceType))

	if newObjs.Len() != 1 {
		return nil, errors.New("expected one object in POST")
	}

	newObj := newObjs.Index(0).Interface()

	// Create model
	id, err := r.source.Create(newObj, req)
	if err != nil {
		return nil, err
	}

	// Find model
	obj, err := r.source.FindOne(id, req)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (r *Resource) HandleUpdate(req *Req, mod RequestModel) (interface{}, error) {
	// Find one model to update
	obj, err := r.source.FindOne(req.Params.Get(":id"), req)
	if err != nil {
		return nil, err
	}

	v, err := mod.Unmarshal(req)
	if err != nil {
		return nil, err
	}

	updatingObjs := reflect.ValueOf(v).Convert(reflect.SliceOf(r.resourceType))

	if updatingObjs.Len() != 1 {
		return nil, errors.New("expected one object in PUT")
	}

	// Update model
	updatingObj := updatingObjs.Index(0).Interface()

	if err := r.source.Update(updatingObj, req); err != nil {
		return nil, err
	}

	obj, err = r.source.FindOne(req.Params.Get(":id"), req)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (r *Resource) HandleDelete(req *Req) (interface{}, error) {
	obj, err := r.source.FindOne(req.Params.Get(":id"), req)
	if err != nil {
		return nil, err
	}

	err = r.source.Delete(req.Params.Get(":id"), req)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (r *Resource) HandleError(req *Req, err error) {
	handleError(req, err)
}

func HandleError(req *Req, err error) {
	handleError(req, err)
}

func handleError(req *Req, err error) {
	// Convert err into an Error
	apiErr, ok := err.(Error)
	if !ok {
		apiErr = WrapErr(err, 500)
	}

	req.Response.WriteHeader(apiErr.HTTPStatus())
	fmt.Fprintln(req.Response, apiErr.HTTPBody())
}
