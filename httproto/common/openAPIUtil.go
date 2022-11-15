package common

import (
	_ "embed"
	"github.com/getkin/kin-openapi/openapi3"
)

type AigesOPenAPI struct {
	Version      string
	Title        string
	APIBase      string
	SvcId        string
	Url          string
	InputSchema  *openapi3.Schema
	OutputSchema *openapi3.Schema
}

func (o *AigesOPenAPI) GenerateOpenAPIJson() (*openapi3.T, error) {
	var a openapi3.T
	a.Info = &openapi3.Info{
		Title:       o.Title,
		Description: "this is the doc openapi3.0 for the aiges service ...",
		Version:     o.Version,
		License: &openapi3.License{
			Name: "Apache2",
			URL:  "https://www.apache.org/licenses/LICENSE-2.0",
		},
	}
	a.Servers = openapi3.Servers{
		&openapi3.Server{
			URL:         o.APIBase,
			Description: "api for the aiges service",
		},
	}
	a.Paths = openapi3.Paths{
		"/": &openapi3.PathItem{
			Post: &openapi3.Operation{
				OperationID: "callAIAbility",
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Value: &openapi3.Response{
							Content: openapi3.Content{
								"application/json": &openapi3.MediaType{
									Schema: &openapi3.SchemaRef{
										Value: o.OutputSchema,
									},
								},
							},
						},
					},
					"default": &openapi3.ResponseRef{
						Value: &openapi3.Response{
							Content: openapi3.Content{
								"application/json": &openapi3.MediaType{
									Schema: &openapi3.SchemaRef{
										Value: o.OutputSchema,
									},
								},
							},
						},
					},
				},
				RequestBody: &openapi3.RequestBodyRef{
					Value: &openapi3.RequestBody{
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{
									Value: o.InputSchema,
								},
							},
						},
					},
				},
			},
		},
	}
	a.OpenAPI = "3.0.3"
	return &a, nil
}
