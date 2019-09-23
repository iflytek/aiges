package jsonpointer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func Example() {
	var document = []byte(`{ 
    "foo": {
      "bar": {
        "baz": [0,"hello!"]
      }
    }
  }`)

	// unmarshal our document into generic go structs
	parsed := map[string]interface{}{}
	// be sure to handle errors in real-world code!
	json.Unmarshal(document, &parsed)

	// parse a json pointer. Pointers can also be url fragments
	// the following are equivelent pointers:
	// "/foo/bar/baz/1"
	// "#/foo/bar/baz/1"
	// "http://example.com/document.json#/foo/bar/baz/1"
	ptr, _ := Parse("/foo/bar/baz/1")

	// evaluate the pointer against the document
	// evaluation always starts at the root of the document
	got, _ := ptr.Eval(parsed)

	fmt.Println(got)
	// Output: hello!
}

// doc pulled from spec:
var docBytes = []byte(`{
  "foo": ["bar", "baz"],
  "": 0,
  "a/b": 1,
  "c%d": 2,
  "e^f": 3,
  "g|h": 4,
  "i\\j": 5,
  "k\"l": 6,
  " ": 7,
  "m~n": 8
}`)

func TestParse(t *testing.T) {
	cases := []struct {
		raw    string
		parsed string
		err    string
	}{
		{"#/", "/", ""},
		{"#/foo", "/foo", ""},
		{"#/foo/", "/foo/", ""},

		{"://", "", "parse ://: missing protocol scheme"},
		{"#7", "", "non-empty references must begin with a '/' character"},
		{"", "", ""},
		{"https://example.com#", "", ""},
	}

	for i, c := range cases {
		got, err := Parse(c.raw)
		if !(err == nil && c.err == "" || err != nil && err.Error() == c.err) {
			t.Errorf("case %d error mismatch. expected: '%s', got: '%s'", i, c.err, err)
			continue
		}

		if c.err == "" && got.String() != c.parsed {
			t.Errorf("case %d string output mismatch: expected: '%s', got: '%s'", i, c.parsed, got.String())
			continue
		}
	}
}

func TestEval(t *testing.T) {
	doc := map[string]interface{}{}
	if err := json.Unmarshal(docBytes, &doc); err != nil {
		t.Errorf("error unmarshaling document json: %s", err.Error())
		return
	}

	cases := []struct {
		ptrstring string
		expect    interface{}
		err       string
	}{
		// "raw" references
		{"", doc, ""},
		{"/foo", doc["foo"], ""},
		{"/foo/0", "bar", ""},
		{"/", float64(0), ""},
		{"/a~1b", float64(1), ""},
		{"/c%d", float64(2), ""},
		{"/e^f", float64(3), ""},
		{"/g|h", float64(4), ""},
		{"/i\\j", float64(5), ""},
		{"/k\"l", float64(6), ""},
		{"/ ", float64(7), ""},
		{"/m~0n", float64(8), ""},
		//
		{"/undefined", nil, ""},

		// url fragment references
		{"#", doc, ""},
		{"#/foo", doc["foo"], ""},
		{"#/foo/0", "bar", ""},
		{"#/", float64(0), ""},
		{"#/a~1b", float64(1), ""},
		{"#/c%25d", float64(2), ""},
		{"#/e%5Ef", float64(3), ""},
		{"#/g%7Ch", float64(4), ""},
		{"#/i%5Cj", float64(5), ""},
		{"#/k%22l", float64(6), ""},
		{"#/%20", float64(7), ""},
		{"#/m~0n", float64(8), ""},

		{"https://example.com#/m~0n", float64(8), ""},

		// bad references
		{"/foo/bar", nil, "invalid array index: bar"},
		{"/foo/3", nil, "index 3 exceeds array length of 2"},
		{"/bar/baz", nil, "invalid JSON pointer: /bar/baz"},
	}

	for i, c := range cases {
		ptr, err := Parse(c.ptrstring)
		if err != nil {
			t.Errorf("case %d unexpected parse error: %s", i, err.Error())
			continue
		}

		got, err := ptr.Eval(doc)
		if !(err == nil && c.err == "" || err != nil && err.Error() == c.err) {
			t.Errorf("case %d error mismatch. expected: '%s', got: '%s'", i, c.err, err)
			continue
		}

		if !reflect.DeepEqual(c.expect, got) {
			t.Errorf("case %d result mismatch. expected: %v, got: %v", i, c.expect, got)
			continue
		}
	}
}

func TestDescendent(t *testing.T) {
	cases := []struct {
		parent string
		path   string
		parsed string
		err    string
	}{
		{"#/", "0", "/0", ""},
		{"/0", "0", "/0/0", ""},
		{"/foo", "0", "/foo/0", ""},
		{"/foo", "0", "/foo/0", ""},
		{"/foo/0", "0", "/foo/0/0", ""},
	}

	for i, c := range cases {
		p, err := Parse(c.parent)
		if err != nil {
			t.Errorf("case %d error parsing parent: %s", i, err.Error())
			continue
		}

		desc, err := p.Descendant(c.path)
		if !(err == nil && c.err == "" || err != nil && err.Error() == c.err) {
			t.Errorf("case %d error mismatch. expected: %s, got: %s", i, c.err, err)
			continue
		}

		if desc.String() != c.parsed {
			t.Errorf("case %d: expected: %s, got: %s", i, c.parsed, desc.String())
			continue
		}
	}
}

func BenchmarkEval(b *testing.B) {
	document := []byte(`{ 
    "foo": {
      "bar": {
        "baz": [0,"hello!"]
      }
    }
  }`)

	parsed := map[string]interface{}{}
	json.Unmarshal(document, &parsed)
	ptr, _ := Parse("/foo/bar/baz/1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := ptr.Eval(parsed); err != nil {
			b.Errorf("error evaluating: %s", err.Error())
			continue
		}

	}
}

func TestEscapeToken(t *testing.T) {
	cases := []struct {
		input  string
		output string
	}{
		{"/abc~1/~/0/~0/", "/abc~1/~/0/~0/"},
	}
	for i, c := range cases {
		got := unescapeToken(escapeToken(c.input))
		if got != c.output {
			t.Errorf("case %d result mismatch.  expected: '%s', got: '%s'", i, c.output, got)
		}
	}
}
