package jsonschema

import "github.com/tidwall/gjson"

type GValidator interface {
	GValidate(ctx *ValidateCtx, val *gjson.Result)
}
