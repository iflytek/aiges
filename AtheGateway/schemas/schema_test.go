package schemas

import (
	"testing"
	"github.com/qri-io/jsonschema"
	"encoding/json"
	"fmt"
	"unsafe"
	"go/types"
	"reflect"
	"time"
	"github.com/oliveagle/jsonpath"
	"log"
	"github.com/json-iterator/go"
)

func TestSchema(t *testing.T) {
	schemaData := []byte(`{
      "title": "Person",
      "type": "object",
      "properties": {
          "firstName": {
              "type": "string",
"magic":{
	"key":"age",
	"enable":true
}
          },
          "lastName": {
              "type": "string"
          },
          "age": {
              "description": "Age in years",
              "type": "integer",
              "minimum": 0
          },
          "friends": {
            "type" : "array",
            "items" : { "title" : "REFERENCE", "$ref" : "#" }
          }
      }
    }`)

	rs := &jsonschema.RootSchema{}
	if err := json.Unmarshal(schemaData, rs); err != nil {
		panic("unmarshal schema: " + err.Error())
	}
	var valid = []byte(`{
    "firstName" : "George",
    "lastName" : "Michael"
    }`)
	var i interface{}
	json.Unmarshal(valid,&i)
	var err  = &[]jsonschema.ValError{}
	if rs.Validate("/",i,err); len(*err) > 0 {
		log.Println(*err)
	}

	log.Println(i)
	time.Now()
}

type Req struct {
	Common struct{
		Appid string `json:"app_id"`
	} `json:"common"`
	Data struct{
		Status int `json:"status"`
	} `json:"data"`
}

func TestValidate(t *testing.T) {
	m:=map[string]interface{}{}
	m["kk"]=12
	ts(&m)
	log.Println(m)
}

func ts(v interface{})  {
	if m,ok:=v.(*map[string]interface{});ok{
		(*m)["kk"]=nil
		log.Println("cl")
	}
	switch v.(type) {

	}
}


func TestMarshalInterface(t *testing.T) {
	var i interface{}
	fmt.Println(marshalInterface("$[0].data",&i,"hello"))
	fmt.Println(marshalInterface("$[0].data_type",&i,1))
	marshalInterface("$[0].encoding",&i,"raw")

	fmt.Println(marshalInterface("$[1].data",&i,"hello_world"))
	fmt.Println(marshalInterface("$[1].data_type",&i,2))
	fmt.Println(marshalInterface("$[1].encoding",&i,"speex"))

	fmt.Println(marshalInterface("$[2].data",&i,"hello_world2"))
	fmt.Println(marshalInterface("$[2].data_type",&i,3))
	fmt.Println(marshalInterface("$[2].encoding",&i,"speex-wb"))
	fmt.Println(marshalInterface("$[3]",&i,"speex-wb"))

	b,_:=json.Marshal(i)
	fmt.Println(string(b))
	fmt.Println(reflect.TypeOf(&i))
}

func TestMarshalInterface2(t *testing.T) {
	var i interface{}
	fmt.Println(marshalInterface("$.data",&i,"hello"))
	fmt.Println(marshalInterface("$.data_type",&i,1))
//	marshalInterface("$[0].encoding",&i,"raw")

	fmt.Println(marshalInterface("$.data1[0]",&i,"hello_world"))
	fmt.Println(marshalInterface("$.data_type1[0]",&i,2))
	fmt.Println(marshalInterface("$.data_type1[3]",&i,3))
	fmt.Println(marshalInterface("$.encoding1[0]",&i,"speex"))
	fmt.Println(marshalInterface("$.desc_args.encoding",&i,"speex"))
	fmt.Println(marshalInterface("$.desc_args.format",&i,"16000"))

	b,_:=json.Marshal(i)
	fmt.Println(string(b))
	fmt.Println(reflect.TypeOf(&i))
}
var jsonss = jsoniter.ConfigCompatibleWithStandardLibrary
func TestJsonnit(t *testing.T) {
	var schemaData = []byte(`{
      "title": "Person",
      "type": "object",
      "properties": {
          "firstName": {
              "type": "string"
          },
          "lastName": {
              "type": "string"
          },
          "age": {
              "description": "Age in years",
              "type": "integer",
              "minimum": 0
          },
          "friends": {
            "type" : "array",
            "items" : { "title" : "REFERENCE", "$ref" : "#" }
          }
      },
      "required": ["firstName", "lastName"]
    }`)

	for k:=0;k<500000;k++{
		var i interface{}
		jsonss.Unmarshal(schemaData,&i)
	}
}

func TestStarnd(t *testing.T) {
	var schemaData = []byte(`{
      "title": "Person",
      "type": "object",
      "properties": {
          "firstName": {
              "type": "string"
          },
          "lastName": {
              "type": "string"
          },
          "age": {
              "description": "Age in years",
              "type": "integer",
              "minimum": 0
          },
          "friends": {
            "type" : "array",
            "items" : { "title" : "REFERENCE", "$ref" : "#" }
          }
      },
      "required": ["firstName", "lastName"]
    }`)

	for k:=0;k<500000;k++{
		var i interface{}
		json.Unmarshal(schemaData,&i)
	}

}
type Data struct {
	Status int
	Data string
}

func TestGetOp(t *testing.T) {
	var i interface{}
	marshalInterface("$[0].ppp",&i,"hello")
	marshalInterface("$[1].ppp",&i,"thk")
	marshalInterface("$[1].kfc",&i,"thk")
	mp:=i.([]interface{})
	fpt:=(*[]interface{})(unsafe.Pointer(&mp))
	fmt.Println(*fpt)
}
var c string
func TestValidateOfMsg(t *testing.T) {
	var s interface{ }= "hello thankyou"
	for i:=0;i<10000000;i++{
		c =fmt.Sprintf("%v",s)
	}

}

func TestValidateOfMsg1(t *testing.T) {
	var s interface{} =  "hello thankyou"
	for i:=0;i<10000000;i++{
		c = (s.(string))
	}

}

type Vad interface {
	Validate()
}

type Tes struct {
	Vad Vad `json:"vad"`
}


func TestProperties_Validate(t *testing.T) {
	var a = Tes{}
	json.Unmarshal([]byte(`{"vad":"sdf"}`),&a)
	fmt.Println(a)
}

func svc(v Vad)  {
	if v==nil{
		return
	}
	v.Validate()
}