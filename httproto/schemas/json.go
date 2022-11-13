package schemas

import (
	"encoding/json"
	"fmt"
)

type JsonElement struct {
	data interface{}
}

func NewJsonElem(i interface{}) *JsonElement {
	return &JsonElement{data: i}
}

func (e *JsonElement) Get(key string) *JsonElement {
	if e == nil {
		return nil
	}
	m, ok := e.data.(map[string]interface{})
	if ok {
		return &JsonElement{
			data: m[key],
		}
	}
	return nil
}

func (e *JsonElement) GetAsString() (string, bool) {
	if e == nil {
		return "", false
	}
	s, ok := e.data.(string)
	return s, ok
}

func (e *JsonElement) Interface() interface{} {
	if e == nil {
		return nil
	}
	return e.data
}

func (e *JsonElement) GetIndex(idx int) *JsonElement {
	if e == nil {
		return nil
	}
	arr, ok := e.data.([]interface{})
	if ok && idx < len(arr) {
		return &JsonElement{data: arr[idx]}
	}
	return nil
}

func (e *JsonElement) Set(key string, val interface{}) error {
	if e == nil {
		return fmt.Errorf("element of %s is nil", key)
	}
	m, ok := e.data.(map[string]interface{})
	if ok {
		m[key] = val
		return nil
	}
	return fmt.Errorf("data of %s is not object", key)
}

func (e *JsonElement) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &e.data)
}

func (e *JsonElement) MarshalJSON() ([]byte, error) {
	if e == nil {
		return nil, nil
	}
	b, err := json.Marshal(e.data)
	return b, err
}

// docker-exec hub.iflytek.com/aiaas/personal-test:latest ./test -ip a.b.c
