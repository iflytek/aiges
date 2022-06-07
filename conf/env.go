package conf

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type envmap struct {
	env   string
	value interface{}
}

var envMaps = []envmap{
	{"CATCH_WRAPPER_DELAY_PERIOD", &WrapperDelayDetectPeriod}, //接口响应超时
}

func parseLoaderEnv() (err error) {
	// 环境变量解析
	for _, v := range envMaps {
		if val := os.Getenv(v.env); len(val) != 0 {
			switch v.value.(type) {
			case *bool:
				var ev int
				ev, err = strconv.Atoi(val)
				if err != nil {
					return errors.New(fmt.Sprintf("env invalid with %s, name:%s, value:%s", err.Error(), v.env, val))
				}
				// Bool环境变量取值：0,1
				switch ev {
				case 0:
					*(v.value.(*bool)) = false
				default:
					*(v.value.(*bool)) = true
				}
			case *string:
				*(v.value.(*string)) = val
			case *int:
				var ev int
				ev, err = strconv.Atoi(val)
				if err != nil {
					return errors.New(fmt.Sprintf("env invalid with %s, name:%s, value:%s", err.Error(), v.env, val))
				}
				*(v.value.(*int)) = ev
			case *uint32:
				var ev int
				ev, err = strconv.Atoi(val)
				if err != nil {
					return errors.New(fmt.Sprintf("env invalid with %s, name:%s, value:%s", err.Error(), v.env, val))
				}
				*(v.value.(*uint32)) = uint32(ev)
			case *uint64:
				var ev int
				ev, err = strconv.Atoi(val)
				if err != nil {
					return errors.New(fmt.Sprintf("env invalid with %s, name:%s, value:%s", err.Error(), v.env, val))
				}
				*(v.value.(*uint64)) = uint64(ev)
			case *[]string:
				// []string环境变量以","作为分隔符
				tmp := strings.Split(val, ",")
				*(v.value.(*[]string)) = nil // clear
				*(v.value.(*[]string)) = append(*(v.value.(*[]string)), tmp...)
			default:
				return errors.New(fmt.Sprintf("env invalid type, name:%s, value:%s", v.env, val))
			}
		}
	}
	return
}
