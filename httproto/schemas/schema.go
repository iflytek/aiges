package schemas

import (
	"encoding/json"
	"log"
	"sync"
)

type SvcSchema struct {
	Content string
}

var once sync.Once
var instance AISchema

func SetSchemaFromPython(s string) *AISchema {
	once.Do(func() {

		err := json.Unmarshal([]byte(s), &instance)
		if err != nil {
			log.Fatalf("wrong schema format... ...check...%s\n", err.Error())
		}
	})
	return &instance
}

func GetSvcSchemaFromPython() *AISchema {
	return &instance
}
