package utils

import (
	"encoding/json"
	"fmt"
)

func AddExtraTag(base string, kvs map[string]string) (extra string, err error) {
	return addExtraTag(base, kvs)
}
func addExtraTag(base string, kvs map[string]string) (extra string, err error) {
	m := make(map[string]string)
	if len(base) != 0 {
		err = json.Unmarshal([]byte(base), &m)
		if err != nil {
			return
		}
	}
	for k, v := range kvs {
		m[k] = v
	}
	jsonRst, jsonRstErr := json.Marshal(m)
	if jsonRstErr != nil {
		err = jsonRstErr
	}
	extra = string(jsonRst)
	return
}
func extractExtraTag(base string, tag string) (val string, err error) {
	m := make(map[string]string)
	if len(base) != 0 {
		err = json.Unmarshal([]byte(base), &m)
		if err != nil {
			return
		}
	}
	valTmp, valTmpOk := m[tag]
	val = valTmp
	if !valTmpOk {
		err = fmt.Errorf("can't extract tag:%v from extra:%v", tag, base)
	}
	return
}
