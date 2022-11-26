package common

import (
	"encoding/json"
)

func replaceSchema(propName string, targetName string, prop interface{}, isPropertiesItems bool) {
	switch e := prop.(type) {
	case map[string]interface{}:
		for key, val := range e {
			if !isPropertiesItems {
				switch key {
				case "properties":
					replaceSchema(propName, targetName, val, true)
				case "if", "else", "then":
					replaceSchema(propName, targetName, val, false)
				default:
					if key == propName {
						delete(e, propName)
						e[targetName] = val
					}

				}
			} else {
				replaceSchema(propName, targetName, val, false)
			}

		}
	case []interface{}:
		for _, v := range e {
			replaceSchema(propName, targetName, v, false)
		}
	}
}

func ReplaceSchema(sc []byte, key string, replace string) ([]byte, error) {
	var i interface{}
	err := json.Unmarshal(sc, &i)
	if err != nil {
		return nil, err
	}
	replaceSchema(key, replace, i, false)
	return json.Marshal(i)
}
