package xsf

import (
	"fmt"
	"time"
)

var globalStart time.Time

func interface2stringslice(rawData interface{}) ([]string, error) {
	interfaceSlice, interfaceSliceOk := rawData.([]interface{})
	if !interfaceSliceOk {
		return nil, fmt.Errorf("can't convert %v to []interface{}", rawData)
	}
	var rstSlice []string
	for _, interfaceSliceItem := range interfaceSlice {
		stringItem, stringItemOk := interfaceSliceItem.(string)
		if !stringItemOk {
			return nil, fmt.Errorf("can't convert %v to string", interfaceSliceItem)
		}
		rstSlice = append(rstSlice, stringItem)
	}
	return rstSlice, nil
}
