package stream

import (
	"github.com/xfyun/aiges/httproto/common"
	"strconv"
	"time"
)

func currentTimeMills() int {
	return int(time.Now().UnixNano()) / int(time.Millisecond)
}

func CodeMapping(codeMap map[string]interface{}, code int) (int, bool) {
	if codeMap == nil {
		return 0, false
	}

	codeStr := strconv.Itoa(code)
	refCode, ok := codeMap[codeStr]
	if !ok {
		return 0, false
	}
	intCode := common.Int(refCode)
	if intCode > 100 && intCode < 600 {
		return intCode, true
	}
	return 0, false
}
