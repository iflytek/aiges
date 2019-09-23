package utils

import (
	"io/ioutil"
	"errors"
	"encoding/json"
)

type xsfcfg struct {
	local  map[string]string
	lb     map[string]string
	log    map[string]string
	common map[string]string
	custom map[string]string
}

func (x *xsfcfg) jsonparser(jsoncfg string) error {
	if "" == jsoncfg {
		return errors.New("json cfg file is empty")
	}
	fc, fe := ioutil.ReadFile(jsoncfg)
	if fe != nil {
		return fe
	}
	var jc struct {
		Framework struct {
			Local  map[string]string `json:"local"`
			Lb     map[string]string `json:"lb"`
			Log    map[string]string `json:"log"`
			Common map[string]string `json:"common"`
		} `json:"framework"`
		Custom map[string]string `json:"custom"`
	}
	if e := json.Unmarshal(fc, &jc); nil != e {
		return e
	}
	x.local = jc.Framework.Local
	x.lb = jc.Framework.Lb
	x.log = jc.Framework.Log
	x.common = jc.Framework.Common
	x.custom = jc.Custom
	return nil
}
