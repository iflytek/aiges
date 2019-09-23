package jsonref

import (
	"encoding/json"
	"github.com/oliveagle/jsonpath"
	"log"
	"testing"
)
var jsons = `
[{
	"bege":{
		"cp":[1,2,3],
		"cps":[
			{
			   "cp":[1,2,3],
				"cp1":[{"haha":"haha"}]
			}
			
		]
	}
}]

`
var i interface{}

func init() {
	json.Unmarshal([]byte(jsons),&i)
}
func TestLookup(t *testing.T) {
	log.Println(jsonpath.JsonPathLookup(i,"$[0].bege"))
	for i:=0;i<1000000 ;i++  {
		there()
	}

}

func TestMy(t *testing.T) {
	log.Println(Lookup("[0].bege",i))
	for i:=0;i<1000000 ;i++  {
		lookup()
	}
}

func lookup()  {
	//var i interface{}
	Lookup("[0].bege.cps[0].cp1[0].haha",i)
}

func there()  {
	jsonpath.JsonPathLookup(i,"$[0].bege.cps[0].cp1[0].haha")
}

func TestLookup2(t *testing.T) {
	log.Println(jsonpath.JsonPathLookup(i,"$[0]"))
	log.Println(jsonpath.JsonPathLookup(i,"$[0].bege.cps[0].cp1[0].haha"))
	log.Println(Lookup("$[0]",i))
	log.Println(Lookup("$[0].bege.cps[0].cp1[0].haha",i))
}

func Test_Token(t *testing.T)  {
	for i:=0;i<1000000;i++{
		tokenize("$[0].bege.cps[0].cp1[0].haha")
	}

}
func Test_Token2(t *testing.T)  {
	for i:=0;i<1000000;i++{
		tokenize2("$[0].bege.cps[0].cp1[0].haha")
	}

}

func TestLookup3(t *testing.T) {
	s:=`{
	"bege":{
		"cp":[1,2,3],
		"cps":[
			{
			   "cp":[1,2,3],
				"cp1":[{"haha":"haha"}]
			}
			
		]
	}
}`

	var i interface{}
	json.Unmarshal([]byte(s),&i)
	log.Println(Lookup("$.bege.sd",i))
	log.Println(Lookup("bege",i))
}

func TestMapp(t *testing.T) {
	var a  = []map[string]interface{}{}
	c(a)
}

func c(i interface{})  {
	log.Println(i.([]interface{}))
}