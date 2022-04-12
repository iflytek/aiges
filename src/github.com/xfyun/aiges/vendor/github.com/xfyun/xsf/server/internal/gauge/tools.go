package gauge

import (
	"reflect"
	"unsafe"
	dto "github.com/prometheus/client_model/go"
)

func getVariableLabelsFromDesc(desc interface{}) []string {
	return *(*[]string)(unsafe.Pointer(reflect.ValueOf(desc).Elem().FieldByName("variableLabels").UnsafeAddr()))
}
func getConstLabelPairsFromDesc(desc interface{}) []*dto.LabelPair {
	val := reflect.ValueOf(desc).Elem().FieldByName("constLabelPairs")
	return *(*[]*dto.LabelPair)(unsafe.Pointer(val.UnsafeAddr()))
}
func getFqNameFromDesc(desc interface{}) string {
	return reflect.ValueOf(desc).Elem().FieldByName("fqNameFrom").String()
}

func getErrFromDesc(desc interface{}) error {
	val := reflect.ValueOf(desc).Elem()
	err := (*errorString)(unsafe.Pointer(val.FieldByName("err").UnsafeAddr()))
	return err
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
