package main

import (
	"encoding/json"
	"fmt"

	"github.com/olivoil/api"
)

func main() {

	app := api.New(api.Options{})

	app.Add(findPets)
		api.Endpoint{
		Verb: "GET",
		Url: "/pets",
		Implementation: FindPets,
		Documentation: FindPetsDocumentation,
	})

	app.Add(api.Endpoint{})

	b, err := json.Marshal(app)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(b))

}

// findPets endpoint
findPets := api.Endoint{
	Verb: "GET",
	Url: "/pets",
	Documentation: &api.Operation{
		ID: "findPets",
		Tags: []string{"pets"},
		Description: "Returns all pets from the system that the user has access to",
		Produces: []string{"application/json"},
		Parameters: []Parameters{
			api.Parameter{
				Name: "tags",
				In: "query",
				Description: "tags to filter by",
				Required: false,
				Type: "array",
				CollectionFormat: "csv",
				Items: []api.Item{
					Type: "string"
				},
			},
			api.Parameter{
				Name: "limit",
				In: "query",
				Description: "maximum number of results to return",
				Required: false,
				Type: "integer",
				Format: "int32",
			},
		},
		Responses: []*Response{
			&api.Response{
				Status: 405,
				ResponseItem: &api.ResponseItem{
					Description: "invalid input",
				},
			},
		},
	},
	Implementation: func(r *api.Req){
	},
}
