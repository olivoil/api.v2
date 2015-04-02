package api

import "fmt"

// Resource to facilitate writing the Implementation function in REST endpoints
// Example Usage:
//
//   type UserDataSource struct {
//     db *sql.DB
//   }
//
//   func (s *UserDataSource) FindOne(model *Model) (*Model, error) {}
//   func (s *UserDataSource) FindAll(model *Model) (*Model, error) {
//     s.db.Limit(model.Query.GetInt(limit)).Offset(model.Query.GetInt(offset)).Find(&model.Data)
//     s.db.Model(&User{}).Count(&model.Response["total"])
//     return model, nil
//   }
//   func (s *UserDataSource) Create(model *Model) (string, error) {}
//   func (s *UserDataSource) Update(model *Model) (string, error) {}
//   func (s *UserDataSource) Delete(model *Model) error {}
//
//   type ListUsersRequest struct {
//   }
//
//   func (r *ListUsersRequest) ParseRequest(req *Req) (*Model, error) {
//     return &Model{
//       Data: []*User{},
//       Query: Meta{
//         "limit": strconv.Atoi(req.Params.Get("limit")),
//         "offset": strconv.Atoi(req.Params.Get("offset")),
//       },
//       Response: Meta{
//         "Location": "/users/"+user.ID,
//       },
//     }, nil
//   }
//
//   type UsersResponse struct {
//     Data []*User `json:"data"`
//     Links map[string]string `json:"links,omitempty"`
//   }
//
//   func (r *UsersResponse) Body(model *Model) interface{} {
//     r.Data = model.Data
//     r.Links["self"] = model.Response.GetString("Location")
//     return r
//   }
//
//   func (r *UsersResponse) Status(model *Model) int {
//      if len(model.Data) > 0 {
//        return  http.StatusOk
//      } else {
//        return http.StatusNotFound
//      }
//   }
//
//   func (r *UsersResponse) Headers(model *Model) map[string]string {
//     return model.Response.GetStringMapString("headers")
//   }
//
//   api.Add(Endpoint{
//     Verb: "GET",
//     Path: "/users",
//     Implementation: func(req *Req){
//       res := NewResource(req, &UserDataSource{db: db})
//
//       model, err := res.HandleIndex(&ListUsersRequest{})
//       if err != nil {
//         res.HandleError(err)
//         return
//       }
//
//       err = res.Send(model, &UsersResponse{})
//       if err != nil {
//         res.HandleError(err)
//         return
//       }
//     },
//   })
//
//   api.Activate(router)
//
type Resource struct {
	Req    *Req
	Source DataSource
}

// DataSource provides methods needed for CRUD.
type DataSource interface {
	// FindOne returns a model from a parsed query
	FindOne(*Model) error

	// FindAll returns all objects specified in query
	FindAll(*Model) error

	// Create a new object and return its ID
	// CONVENTION: place created resource id in model.Response["id"]
	Create(*Model) error

	// Update an object and return its ID
	Update(*Model) error

	// Delete an object
	Delete(*Model) error
}

func NewResource(req *Req, source DataSource) *Resource {
	return &Resource{
		Source: source,
		Req:    req,
	}
}

func (r *Resource) HandleIndex(rp RequestParser) (model *Model, err error) {
	model, err = rp.ParseRequest(r.Req)
	if err != nil {
		return
	}

	err = r.Source.FindAll(model)
	return
}

func (r *Resource) HandleRead(rp RequestParser) (model *Model, err error) {
	model, err = rp.ParseRequest(r.Req)
	if err != nil {
		return
	}

	err = r.Source.FindOne(model)
	return
}

func (r *Resource) HandleCreate(rp RequestParser) (model *Model, err error) {
	// Unmarshal request model into model values
	model, err = rp.ParseRequest(r.Req)
	if err != nil {
		return
	}

	err = r.Source.Create(model)
	if err != nil {
		return
	}

	model.Query.Set("id", model.Response.Get("id"))
	err = r.Source.FindOne(model)
	return
}

func (r *Resource) HandleUpdate(rp RequestParser) (model *Model, err error) {
	// Unmarshal request model into model values
	// CONVENTION: pr.Data => pointer to model struct containing the properties that will be updated
	model, err = rp.ParseRequest(r.Req)
	if err != nil {
		return
	}

	err = r.Source.Update(model)
	if err != nil {
		return
	}

	err = r.Source.FindOne(model)
	return
}

func (r *Resource) HandleDelete(rp RequestParser) (model *Model, err error) {
	model, err = rp.ParseRequest(r.Req)
	if err != nil {
		return
	}

	err = r.Source.FindOne(model)
	if err != nil {
		return
	}

	err = r.Source.Delete(model)
	return
}

func (r *Resource) HandleError(err error) {
	handleError(r.Req, err)
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

func (r *Resource) Send(model *Model, rm ResponseMarshaller) error {
	return send(r.Req, model, rm)
}

func Send(req *Req, model *Model, rm ResponseMarshaller) error {
	return send(req, model, rm)
}

func send(req *Req, model *Model, rm ResponseMarshaller) error {
	encoder := JsonEncoder{}
	data, err := encoder.Encode(rm.Body(model))
	if err != nil {
		return err
	}
	headers := rm.Headers(model)
	if headers != nil {
		for key, val := range headers {
			req.Response.Header().Set(key, val)
		}
	}
	req.Response.WriteHeader(rm.Status(model))
	req.Response.Write(data)
	return nil
}
