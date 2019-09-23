package utils

import "testing"

func TestJsonparser(t *testing.T) {
	var xsfparser xsfcfg
	if e := xsfparser.jsonparser("xsf.json"); e != nil {
		t.Fatal(e)
	}
	t.Logf("xsf:%+v\n", xsfparser)
}
// xsf.json
// {
//   "framework": {
//     "local": {
//       "retryTimes": "xxx",
//       "timeout": "xxx"
//     },
//     "lb": {
//       "lbAddr": "xxx",
//       "preCon": "xxx"
//     },
//     "log": {
//       "file": "xxx",
//       "level": "xxx",
//       "maxBackups": "xxx",
//       "maxSize": "xxx",
//       "maxAge": "xxx"
//     },
//     "common": {}
//   },
//   "custom": {
//     "k1": "v1",
//     "k2": "v2",
//     "k3": "v3"
//   }
// }