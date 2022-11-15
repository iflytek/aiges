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

func SetSchema(s string) *AISchema {
	once.Do(func() {

		err := json.Unmarshal([]byte(s), &instance)
		if err != nil {
			log.Fatal("wrong schema format... ...check...")
		}
	})
	return &instance
}

func GetSvcSchema() *AISchema {
	return &instance
}
