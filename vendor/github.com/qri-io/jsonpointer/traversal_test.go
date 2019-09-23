package jsonpointer

import (
	"encoding/json"
	"testing"
)

type A struct {
	Foo  string `json:"$foo,omitempty"`
	Foo2 bool   `json:"$foo2"`
	Bar  B
}

type B struct {
	nope string
	Baz  int
	C    C
	D    D
}

type C struct {
	Nope string
}

func (c C) JSONProps() map[string]interface{} {
	return map[string]interface{}{
		"bat":   "book",
		"stuff": false,
		"other": nil,
	}
}

type D []string

var data = []byte(`{
  "$foo" : "fooval",
  "$foo2" : true,
  "bar" : {
    "baz" : 1,
    "C" : {
      "won't" : "register"
    }
  }
}`)

func TestWalkJSON(t *testing.T) {
	a := &A{}
	if err := json.Unmarshal(data, a); err != nil {
		t.Errorf("unexpected unmarshal error: %s", err.Error())
		return
	}

	elements := 0
	expectElements := 9
	WalkJSON(a, func(elem interface{}) error {
		t.Logf("%#v", elem)
		elements++
		return nil
	})

	if elements != expectElements {
		t.Errorf("expected %d elements, got: %d", expectElements, elements)
	}

}
