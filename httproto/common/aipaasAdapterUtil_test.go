package common

import (
	"fmt"
	"testing"
)

func TestResc(t *testing.T) {
	sc := `
{
 "type":"object",
 "properties":{
  "default":{
   "type":"string",
   "default":"5"
  }
 },
 "if":{
  "type":"properties"
 },
 "anyOf":[
  {
   "type":"object",
   "properties":{
    "default":{
     "type":"string",
     "default":"5"
    }
   }
  }
 ]
 
}

`
	res, err := ReplaceSchema([]byte(sc), "default", "defaultVal")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))
}
