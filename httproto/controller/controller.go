package controller

import (
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/xfyun/aiges/httproto/common"
	"github.com/xfyun/aiges/httproto/schemas"
	"log"
	"net/http"
	"strings"
)

func GetOpenAPIJSON(c *gin.Context) {
	svc := schemas.GetSvcSchema()
	inputSchemaJson, err := svc.InputSchema.MarshalJSON()
	if err != nil {
		log.Println(err.Error())
		return
	}
	outputSchemaJson, err := svc.SchemaOutput.MarshalJSON()
	if err != nil {
		log.Println(err.Error())
		return
	}
	var inputOpenAPI openapi3.Schema
	var outputOpenAPI openapi3.Schema
	json.Unmarshal(inputSchemaJson, &inputOpenAPI)
	json.Unmarshal(outputSchemaJson, &outputOpenAPI)

	var api = common.AigesOPenAPI{
		APIBase:      strings.Join(svc.Meta.GetRoute(), ""), // ??
		SvcId:        svc.Meta.GetServiceId(),
		Title:        fmt.Sprintf("Aiges Service doc for %s", svc.Meta.GetServiceId()),
		InputSchema:  &inputOpenAPI,
		OutputSchema: &outputOpenAPI,
	}
	s, _ := api.GenerateOpenAPIJson()
	c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.IndentedJSON(http.StatusOK, s)

}
