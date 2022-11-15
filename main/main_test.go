package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
)

//go:embed test_schema.json
var tests []byte

func main_test() {
	var a openapi3.T
	a.Info = &openapi3.Info{
		Title:       "new test",
		Description: "descccc",
		Version:     "1.1.1",
		License: &openapi3.License{
			Name: "Apache2",
			URL:  "https://www.apache.org/licenses/LICENSE-2.0",
		},
	}
	a.Servers = openapi3.Servers{
		&openapi3.Server{
			URL:         "/api/v3",
			Description: "hello",
		},
	}
	var ins openapi3.Schema
	err := json.Unmarshal(tests, &ins)
	if err != nil {
		fmt.Println(err)
	}
	a.Paths = openapi3.Paths{
		"scccc": &openapi3.PathItem{
			Post: &openapi3.Operation{
				OperationID: "first",
				RequestBody: &openapi3.RequestBodyRef{
					Value: &openapi3.RequestBody{
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{
									Value: &ins,
								},
							},
						},
					},
				},
			},
		},
	}
	a.OpenAPI = "3.0.2"
	out, err := a.MarshalJSON()
	fmt.Println(string(out))
}
