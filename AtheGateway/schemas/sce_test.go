package schemas

import (
	"encoding/json"
	"fmt"
	"github.com/qri-io/jsonschema"
	"testing"
	"github.com/skyniu/routine-local"
	"runtime"
	"sync"
	"regexp"
)

type As int
type Ints interface {
	Show()
}

func pac()  {
	defer func() {
		if err:=recover();err !=nil{

		}
	}()
	panic("haha")
}
func TestGetByJPath(t *testing.T) {
	for i:=0;i<10000000;i++{

		pac()
	}
}
func aa(i interface{})  {
	s,ok:=i.(*int)
	fmt.Println(s,ok)
}

func main1() {
	//注册校验器
	jsonschema.RegisterValidator("properties", NewProperties)

	schemaData := []byte(`{
      "title": "Person",
      "type": "object",
      "properties": {
          "firstName": {
              "type": "string",
				"$id":"#w"
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
          },
			"name":{
				"$ref":"#w"
			}
			
      }
    }`)
	rs := &jsonschema.RootSchema{}
	if err := json.Unmarshal(schemaData, rs); err != nil {
		panic("unmarshal schema: " + err.Error())
	}
	var valid = []byte(`{
	"firstName" : "sd",
	"lastName" : "Michael",
"name":"sd"
	
	}`)
	var doc interface{}
	json.Unmarshal(valid, &doc)
	fmt.Println("validJson===>", doc)
	errs := []jsonschema.ValError{}
	rs.Validate("/", doc, &errs)
	for _, err := range errs {
		fmt.Println("===>", err.Message)
	}

	//if errors, _ := rs.ValidateBytes(valid); len(errors) > 0 {
	//	for _, err := range errors {
	//		fmt.Println(err.Message)
	//	}
	//}

}

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func TestLoadRoteMapping(t *testing.T) {
	for i:=0;i<1000000;i++{
		runtime.Goid()
	}
}

func TestMarshal(t *testing.T) {
	for i:=0;i<1000000;i++{
		rlocal.GetGID()
	}
}
type Meta struct {
	Name string
	Age int
	Mail string
	Male bool
	Yes string
	No string
	Hello int
	Session sync.Mutex

}
var pool = sync.Pool{}

func New() *Meta {
	return &Meta{}
}

func TestPool(t *testing.T) {
	pool.New = func() interface{} {
		return &Meta{}
	}
	for i:=0;i<1000000;i++{
		//s:=pool.Get()
		//pool.Put(s)
		New()
	}

}


type data struct {
	name string
}
type IT struct {
	Fi string
}



func TestGetMapping(t *testing.T) {
	var rg= regexp.MustCompile(`\$(\.\w+(\[\d+\])?)+$`)
	b:=rg.MatchString("$.abc2[123].acd2[12]")
	fmt.Println(b)
	fmt.Println(rg.MatchString("$.abc2.acd2[1].sdfds.sdf"))

}

func a(v interface{})  {
	s:=v.([]interface{})
	fmt.Println(s)
}

func TestGetRespByCall(t *testing.T) {
	var exprg = regexp.MustCompile(`(.+)\((.+)\)$`)
	fmt.Println(exprg.MatchString("len(data)"))
	r:=exprg.FindAllSubmatch([]byte("len(data)"),-1)
	rg:=regexp.MustCompile(`if\((\w+)(=|>|!=|<)(\w+)\){((\w+)=(\w+);)*}$`)
	fmt.Println(rg.MatchString(`if(a=b){a=b;}`))
	r =rg.FindAllSubmatch([]byte(`if(a=b){a=b;c=d;}`),-1)
	for _,v:=range r{
		//fmt.Println(string(v[0]))
		fmt.Println("----")
		for _,vv:=range v{
			fmt.Println(string(vv))
		}
	}
}

func TestGetRespByCall1(t *testing.T) {
	var exprg = regexp.MustCompile(`(hmac_username|api_key)="(.+)"\s*headers="(.+)"$`)
	fmt.Println(exprg.MatchString("15029032042"))
	r:=exprg.FindAllSubmatch([]byte(`hmac_username="sdfdsf" headers="fsdfsf"`),-1)
	for _,v:=range r{
		//fmt.Println(string(v[0]))
		fmt.Println("----")
		for _,vv:=range v{
			fmt.Println(string(vv))
		}
	}


}

