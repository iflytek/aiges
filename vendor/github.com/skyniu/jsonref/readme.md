#### josnref for go

**usage**

 Jsonref is a tool used for creating json reversed. with this tool ,you can create a json string by specific query string like "$.user.name" 
 and set a value of "bob" to it. The tool will generate a json string such as follow:
 ```text
{
    "user":{
        "name":"bob"
    }
}
```  

**get start**

- install
```text
    go get -t github.com/skyniu/jsonref
``` 

- use 
```text
import (
	"fmt"
	"encoding/json"
	"github.com/skyniu/jsonref"
)
type User struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

func main() {
	var m = map[string]interface{}{}
	jsonref.Marshal("$.zhansan",m,User{"zhangsan",20})
	jsonref.Marshal("$.class[0]",m,User{"lisi",11})
	jsonref.Marshal("$.class[1]",m,User{"wangwu",18})
	jsonref.Marshal("$.class[2]",m,User{"dajj",18})
	jsonref.Marshal("$.group[5].name",m,"biaoge")
	jsonref.Marshal("$.group[5].age",m,12)

	s,_:=json.Marshal(m)
	fmt.Println(string(s))
}

```
the output json is :
```text
{
  "class": [
    {
      "name": "lisi",
      "age": 11
    },
    {
      "name": "wangwu",
      "age": 18
    },
    {
      "name": "dajj",
      "age": 18
    }
  ],
  "group": [
    {},
    {},
    {},
    {},
    {},
    {
      "age": 12,
      "name": "biaoge"
    }
  ],
  "zhansan": {
    "name": "zhangsan",
    "age": 20
  }
}
```