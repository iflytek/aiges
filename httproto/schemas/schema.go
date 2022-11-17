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
			log.Fatal("wrong schema format... ...check...")
		}
	})
	return &instance
}

func GetSvcSchemaFromPython() *AISchema {
	return &instance
}
